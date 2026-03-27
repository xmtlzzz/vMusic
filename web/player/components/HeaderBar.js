import { computed, defineComponent } from "https://unpkg.com/vue@3/dist/vue.esm-browser.prod.js";

export const HeaderBar = defineComponent({
  name: "HeaderBar",
  props: {
    session: { type: Object, required: true },
    totalTracks: { type: Number, default: 0 },
    likedCount: { type: Number, default: 0 },
    currentTrack: { type: Object, default: null },
  },
  emits: ["logout"],
  setup(props) {
    const statusText = computed(() => {
      if (!props.session.user) {
        return "登录后可以同步你的我喜欢列表与账号状态。";
      }
      return `欢迎回来，${props.session.user.username}`;
    });

    return { statusText };
  },
  template: `
    <header class="topbar glass-panel">
      <div class="topbar-brand">
        <span class="brand-dot"></span>
        <div>
          <p class="eyebrow">vMusic Web Player</p>
          <h1>云感玻璃播放器</h1>
          <p class="topbar-copy">{{ statusText }}</p>
        </div>
      </div>

      <div class="topbar-side">
        <div class="header-stats">
          <div class="header-chip">
            <span class="chip-label">本地曲目</span>
            <strong>{{ totalTracks }}</strong>
          </div>
          <div class="header-chip">
            <span class="chip-label">我喜欢</span>
            <strong>{{ likedCount }}</strong>
          </div>
          <div class="header-chip header-chip-wide">
            <span class="chip-label">正在播放</span>
            <strong>{{ currentTrack ? currentTrack.music_name : '等待选择曲目' }}</strong>
          </div>
        </div>

        <div class="session-actions">
          <span class="status-pill" :class="{ active: !!session.user }">{{ session.user ? '已登录' : '游客模式' }}</span>
          <button v-if="session.user" class="ghost-button compact" type="button" @click="$emit('logout')">退出登录</button>
        </div>
      </div>
    </header>
  `,
});
