import { defineComponent } from 'vue';

export default defineComponent({
  props: {
    englishName: {
      type: String,
      default: '',
    },
    chineseName: {
      type: String,
      default: '',
    },
  },
  setup(props) {
    return () => (
      <div class="bk-metrics-staff-li flex-row align-items-center flex-1 ph10">
        <img
          class="bk-metrics-staff-memeber-pic"
          src={`//dayu.woa.com/avatars/${props.englishName}/profile.jpg`}
        />
        {props.englishName} ({props.chineseName})
      </div>
    );
  },
});
