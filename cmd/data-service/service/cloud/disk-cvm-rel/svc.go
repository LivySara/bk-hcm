/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 混合云管理平台 (BlueKing - Hybrid Cloud Management System) available.
 * Copyright (C) 2022 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * We undertake not to change the open source license (MIT license) applicable
 *
 * to the current version of the project delivered to anyone in the future.
 */

package diskcvmrel

import (
	"fmt"

	"hcm/pkg/api/core"
	datarelproto "hcm/pkg/api/data-service/cloud"
	dataproto "hcm/pkg/api/data-service/cloud/disk"
	"hcm/pkg/criteria/enumor"
	"hcm/pkg/criteria/errf"
	"hcm/pkg/dal/dao"
	reldao "hcm/pkg/dal/dao/cloud/disk-cvm-rel"
	"hcm/pkg/dal/dao/orm"
	"hcm/pkg/dal/dao/tools"
	"hcm/pkg/dal/dao/types"
	tablecloud "hcm/pkg/dal/table/cloud"
	"hcm/pkg/rest"

	"github.com/jmoiron/sqlx"
)

type relSvc struct {
	dao.Set
	objectDao *reldao.DiskCvmRelDao
}

// Init ...
func (svc *relSvc) Init() {
	d := &reldao.DiskCvmRelDao{}
	registeredDao := svc.GetObjectDao(d.Name())
	if registeredDao == nil {
		d.ObjectDaoManager = new(dao.ObjectDaoManager)
		svc.RegisterObjectDao(d)
	}

	svc.objectDao = svc.GetObjectDao(d.Name()).(*reldao.DiskCvmRelDao)
}

// BatchCreate ...
func (svc *relSvc) BatchCreate(cts *rest.Contexts) (interface{}, error) {
	req := new(datarelproto.DiskCvmRelBatchCreateReq)
	if err := cts.DecodeInto(req); err != nil {
		return nil, errf.NewFromErr(errf.DecodeRequestFailed, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	_, err := svc.Txn().AutoTxn(cts.Kit, func(txn *sqlx.Tx, opt *orm.TxnOption) (interface{}, error) {
		rels := make([]*tablecloud.DiskCvmRelModel, len(req.Rels))
		for idx, relReq := range req.Rels {
			rels[idx] = &tablecloud.DiskCvmRelModel{
				CvmID:   relReq.CvmID,
				DiskID:  relReq.DiskID,
				Creator: cts.Kit.User,
			}
		}

		return nil, svc.objectDao.BatchCreateWithTx(cts.Kit, txn, rels)
	})

	return nil, err
}

// List ...
func (svc *relSvc) List(cts *rest.Contexts) (interface{}, error) {
	req := new(datarelproto.DiskCvmRelListReq)
	if err := cts.DecodeInto(req); err != nil {
		return nil, err
	}

	if err := req.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	opt := &types.ListOption{
		Fields: req.Fields,
		Filter: req.Filter,
		Page:   req.Page,
	}
	data, err := svc.objectDao.List(cts.Kit, opt)
	if err != nil {
		return nil, fmt.Errorf("list disk cvm rels failed, err: %v", err)
	}

	if req.Page.Count {
		return &datarelproto.DiskCvmRelListResult{Count: data.Count}, nil
	}

	details := make([]*datarelproto.DiskCvmRelResult, len(data.Details))
	for idx, r := range data.Details {
		details[idx] = &datarelproto.DiskCvmRelResult{
			ID:        r.ID,
			DiskID:    r.DiskID,
			CvmID:     r.CvmID,
			Creator:   r.Creator,
			CreatedAt: r.CreatedAt.String(),
		}
	}

	return &datarelproto.DiskCvmRelListResult{Details: details}, nil
}

// BatchDelete ...
func (svc *relSvc) BatchDelete(cts *rest.Contexts) (interface{}, error) {
	req := new(datarelproto.DiskCvmRelDeleteReq)
	if err := cts.DecodeInto(req); err != nil {
		return nil, err
	}

	if err := req.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	opt := &types.ListOption{
		Fields: []string{"id"},
		Filter: req.Filter,
		Page:   core.DefaultBasePage,
	}

	relResult, err := svc.objectDao.List(cts.Kit, opt)
	if err != nil {
		return nil, fmt.Errorf("list disk cvm rels failed, err: %v", err)
	}

	if len(relResult.Details) == 0 {
		return nil, nil
	}

	delIDs := make([]uint64, len(relResult.Details))
	for idx, rel := range relResult.Details {
		delIDs[idx] = rel.ID
	}

	_, err = svc.Txn().AutoTxn(cts.Kit, func(txn *sqlx.Tx, opt *orm.TxnOption) (interface{}, error) {
		return nil, svc.objectDao.DeleteWithTx(cts.Kit, txn, tools.ContainersExpression("id", delIDs))
	})
	return nil, err
}

// ListWithDisk ...
func (svc *relSvc) ListWithDisk(cts *rest.Contexts) (interface{}, error) {
	req := new(datarelproto.DiskCvmRelWithDiskListReq)
	if err := cts.DecodeInto(req); err != nil {
		return nil, errf.NewFromErr(errf.DecodeRequestFailed, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	data, err := svc.objectDao.ListJoinDisk(cts.Kit, req.CvmIDs)
	if err != nil {
		return nil, err
	}

	disks := make([]*datarelproto.DiskWithCvmID, len(data.Details))
	for idx, d := range data.Details {
		disks[idx] = toProtoDiskWithCvmID(d)
	}
	return disks, nil
}

// ListWithDiskExt ...
func (svc *relSvc) ListWithDiskExt(cts *rest.Contexts) (interface{}, error) {
	vendor := enumor.Vendor(cts.Request.PathParameter("vendor"))
	if err := vendor.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	req := new(datarelproto.DiskCvmRelWithDiskExtListReq)
	if err := cts.DecodeInto(req); err != nil {
		return nil, err
	}

	if err := req.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	data, err := svc.objectDao.ListJoinDisk(cts.Kit, req.CvmIDs)
	if err != nil {
		return nil, err
	}

	switch vendor {
	case enumor.TCloud:
		return toProtoDiskExtWithCvmIDs[dataproto.TCloudDiskExtensionResult](data)
	case enumor.Aws:
		return toProtoDiskExtWithCvmIDs[dataproto.AwsDiskExtensionResult](data)
	case enumor.Gcp:
		return toProtoDiskExtWithCvmIDs[dataproto.GcpDiskExtensionResult](data)
	case enumor.Azure:
		return toProtoDiskExtWithCvmIDs[dataproto.AzureDiskExtensionResult](data)
	case enumor.HuaWei:
		return toProtoDiskExtWithCvmIDs[dataproto.HuaWeiDiskExtensionResult](data)
	default:
		return nil, errf.Newf(errf.InvalidParameter, "unsupported vendor: %s", vendor)
	}
}
