<script setup>
import { formatTime } from "../utils/visual";

defineProps({
  currentTrack: { type: Object, default: null },
  player: { type: Object, required: true },
  discStyle: { type: Object, required: true },
  discLabelStyle: { type: Object, required: true },
  ambientStyle: { type: Object, required: true },
  currentTrackLiked: { type: Boolean, default: false },
  queueCount: { type: Number, default: 0 },
  stageFactItems: { type: Array, default: () => [] },
  stageMetaItems: { type: Array, default: () => [] },
});

defineEmits(["toggle-like", "toggle-play", "previous", "next", "seek"]);
</script>

<template>
  <section class="player-stage glass-panel glass-panel-hero">
    <div class="player-stage__ambient" :style="ambientStyle"></div>
    <div class="player-stage__grid">
      <div class="player-stage__left">
        <div class="scene-badge">Inspired by NetEase Cloud Music & QQ Music</div>

        <div class="disc-stage">
          <div class="disc-stage__glow"></div>
          <div class="tonearm" :class="{ playing: player.isPlaying }"></div>
          <div class="disc-shell" :style="discStyle">
            <div class="disc-shell__groove"></div>
            <div class="disc-center" :style="discLabelStyle">
              <span class="disc-kicker">Now Playing</span>
              <strong>{{ currentTrack ? currentTrack.music_name : '等待播放' }}</strong>
              <span>{{ currentTrack ? currentTrack.author : '选择一首歌开始沉浸播放' }}</span>
            </div>
          </div>
        </div>

        <div class="equalizer" :class="{ active: player.isPlaying }">
          <span></span>
          <span></span>
          <span></span>
          <span></span>
          <span></span>
        </div>
      </div>

      <div class="player-stage__right">
        <div class="player-stage__headline">
          <p class="eyebrow">播放舞台</p>
          <h2>{{ currentTrack ? currentTrack.music_name : '上传音乐后即可开始播放' }}</h2>
          <p class="player-subcopy">{{ currentTrack ? currentTrack.author : '这里会展示当前歌曲信息、播放进度、收藏状态和舞台动效。' }}</p>
        </div>

        <div class="scene-chip-list">
          <div v-for="item in stageFactItems" :key="item.label" class="scene-chip">
            <span>{{ item.label }}</span>
            <strong>{{ item.value }}</strong>
          </div>
        </div>

        <div class="progress-card">
          <div class="time-row">
            <span>{{ formatTime(player.currentTime) }}</span>
            <span>{{ formatTime(player.duration) }}</span>
          </div>
          <input
            type="range"
            min="0"
            :max="player.duration || 0"
            :value="player.currentTime || 0"
            step="0.1"
            @input="$emit('seek', Number($event.target.value))"
          >
        </div>

        <div class="control-row control-row--stage">
          <button class="ghost-button" type="button" @click="$emit('previous')">上一首</button>
          <button class="primary-button play-button" type="button" @click="$emit('toggle-play')">{{ player.isPlaying ? '暂停' : '播放' }}</button>
          <button class="ghost-button" type="button" @click="$emit('next')">下一首</button>
          <button class="like-button" :class="{ active: currentTrackLiked }" type="button" @click="$emit('toggle-like')" :disabled="!currentTrack">
            {{ currentTrackLiked ? '已喜欢' : '加入我喜欢' }}
          </button>
        </div>

        <div class="stage-rows">
          <div v-for="row in stageMetaItems" :key="row" class="stage-row">{{ row }}</div>
        </div>

        <div class="stage-footer">
          <span class="stage-footer__count">当前队列 {{ queueCount }} 首</span>
          <span class="stage-footer__file">{{ currentTrack?.file_name || '等待音频输入' }}</span>
        </div>
      </div>
    </div>
  </section>
</template>
