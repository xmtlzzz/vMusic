<script setup>
import { computed } from "vue";

const props = defineProps({
  session: { type: Object, required: true },
  totalTracks: { type: Number, default: 0 },
  likedCount: { type: Number, default: 0 },
  currentTrack: { type: Object, default: null },
});

defineEmits(["logout"]);

const statusText = computed(() => {
  if (!props.session.user) {
    return "登录后可同步你的喜欢列表、上传记录和播放状态。";
  }
  return `欢迎回来，${props.session.user.username}`;
});
</script>

<template>
  <header class="topbar glass-panel">
    <div class="topbar-brand">
      <span class="brand-dot"></span>
      <div>
        <p class="eyebrow">vMusic Glass Player</p>
        <h1>云感播放器</h1>
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
          <span class="chip-label">当前播放</span>
          <strong>{{ currentTrack ? currentTrack.music_name : '等待选择曲目' }}</strong>
        </div>
      </div>

      <div class="session-actions">
        <span class="status-pill" :class="{ active: !!session.user }">{{ session.user ? '已登录' : '游客模式' }}</span>
        <button v-if="session.user" class="ghost-button compact" type="button" @click="$emit('logout')">退出登录</button>
      </div>
    </div>
  </header>
</template>
