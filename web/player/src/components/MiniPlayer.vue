<script setup>
import { computed } from "vue";
import { gradientStyle, makeInitials } from "../utils/visual";

const props = defineProps({
  currentTrack: { type: Object, default: null },
  isPlaying: { type: Boolean, default: false },
  liked: { type: Boolean, default: false },
  progress: { type: Number, default: 0 },
});

defineEmits(["toggle-play", "previous", "next", "toggle-like"]);

const barStyle = computed(() => ({ width: `${props.progress || 0}%` }));
const coverStyle = computed(() => gradientStyle(props.currentTrack));
const initials = computed(() => makeInitials(props.currentTrack?.music_name || props.currentTrack?.file_name || "VM"));
const metaText = computed(() => {
  if (!props.currentTrack) {
    return "本地音乐库";
  }
  return `${String(props.currentTrack.format || "audio").toUpperCase()} · ${props.currentTrack.from || "local"}`;
});
</script>

<template>
  <aside class="mini-player" :class="{ active: !!currentTrack }">
    <div class="mini-player__progress"><span :style="barStyle"></span></div>
    <div class="mini-player__body glass-panel">
      <div class="mini-player__track">
        <div class="mini-player__cover" :style="coverStyle">{{ initials }}</div>
        <div class="mini-player__copy">
          <strong>{{ currentTrack ? currentTrack.music_name : '等待播放' }}</strong>
          <span>{{ currentTrack ? currentTrack.author : '选择一首歌开始播放' }}</span>
        </div>
      </div>

      <div class="mini-player__controls">
        <button class="icon-button" type="button" @click="$emit('previous')">⏮</button>
        <button class="primary-button mini-player__play" type="button" @click="$emit('toggle-play')">{{ isPlaying ? '暂停' : '播放' }}</button>
        <button class="icon-button" type="button" @click="$emit('next')">⏭</button>
      </div>

      <div class="mini-player__meta">
        <button class="icon-button" :class="{ active: liked }" type="button" @click="$emit('toggle-like')">{{ liked ? '♥' : '♡' }}</button>
        <span>{{ metaText }}</span>
      </div>
    </div>
  </aside>
</template>
