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

package cloud

// SecurityGroupCvmRel define security group cvm rel.
type SecurityGroupCvmRel struct {
	ID              uint64 `json:"id"`
	CvmID           string `json:"cvm_id"`
	SecurityGroupID string `json:"security_group_id"`
	Creator         string `json:"creator"`
	CreatedAt       string `json:"created_at"`
}

// SGCvmRelWithBaseSecurityGroup define security group with cvm id.
type SGCvmRelWithBaseSecurityGroup struct {
	BaseSecurityGroup `json:",inline"`
	CvmID             string `json:"cvm_id"`
	RelCreator        string `db:"rel_creator" json:"rel_creator"`
	RelCreatedAt      string `db:"rel_created_at" json:"rel_created_at"`
}
