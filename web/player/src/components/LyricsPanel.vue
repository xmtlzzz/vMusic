<script setup>
defineProps({
  currentTrack: { type: Object, default: null },
  lines: { type: Array, default: () => [] },
  activeIndex: { type: Number, default: 0 },
});
</script>

<template>
  <section class="lyrics-panel glass-panel">
    <div class="panel-head panel-head--stack">
      <div>
        <p class="eyebrow">歌词区域</p>
        <h3>{{ currentTrack ? '动态歌词面板' : '等待歌曲进入歌词区' }}</h3>
      </div>
      <span class="lyrics-tip">当前为动态歌词占位实现，后续可无缝接入真实 LRC</span>
    </div>

    <div class="lyrics-card">
      <p v-if="currentTrack" class="lyrics-now">{{ currentTrack.music_name }} · {{ currentTrack.author }}</p>
      <div class="lyrics-lines">
        <p
          v-for="(line, index) in lines"
          :key="`${line.time}-${line.text}`"
          class="lyrics-line"
          :class="{ active: index === activeIndex }"
        >
          {{ line.text }}
        </p>
      </div>
    </div>
  </section>
</template>
