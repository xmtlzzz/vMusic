import {
  computed,
  nextTick,
  onBeforeUnmount,
  reactive,
  ref,
} from "https://unpkg.com/vue@3/dist/vue.esm-browser.prod.js";

import { ambientStyle, gradientStyle } from "../utils/visual.js";

export function usePlayer(options = {}) {
  const audioRef = ref(null);
  const state = reactive({
    currentTrack: null,
    isPlaying: false,
    currentTime: 0,
    duration: 0,
    rotation: 0,
    rafId: 0,
    lastFrameTime: 0,
  });

  const discStyle = computed(() => ({
    transform: `rotate(${state.rotation}deg)`,
  }));

  const discLabelStyle = computed(() => gradientStyle(state.currentTrack));
  const stageAmbientStyle = computed(() => ambientStyle(state.currentTrack));

  async function playTrack(track, config = {}) {
    if (!track) {
      return;
    }

    const autoplay = Boolean(config.autoplay);
    const audio = audioRef.value;
    if (!audio) {
      return;
    }

    const sameTrack = state.currentTrack && state.currentTrack.file_name === track.file_name;
    if (!sameTrack) {
      state.rotation = 0;
      state.currentTrack = track;
      state.currentTime = 0;
      state.duration = 0;
      stopRotationLoop();
      state.isPlaying = false;
      audio.src = track.audio_url;
      audio.load();
    } else {
      state.currentTrack = track;
    }

    if (autoplay) {
      try {
        await nextTick();
        await audio.play();
      } catch (error) {
        if (typeof options.onAutoplayBlocked === "function") {
          options.onAutoplayBlocked(error);
        }
      }
    }
  }

  async function togglePlay(fallbackTrack) {
    const audio = audioRef.value;
    if (!audio) {
      return;
    }

    if (!state.currentTrack && fallbackTrack) {
      await playTrack(fallbackTrack, { autoplay: true });
      return;
    }

    if (!state.currentTrack) {
      return;
    }

    if (audio.paused) {
      try {
        await audio.play();
      } catch (error) {
        if (typeof options.onAutoplayBlocked === "function") {
          options.onAutoplayBlocked(error);
        }
      }
      return;
    }

    audio.pause();
  }

  function syncTrack(track) {
    if (!track) {
      return;
    }
    if (state.currentTrack && state.currentTrack.file_name === track.file_name) {
      state.currentTrack = track;
    }
  }

  function clearTrack() {
    const audio = audioRef.value;
    stopRotationLoop();
    state.currentTrack = null;
    state.currentTime = 0;
    state.duration = 0;
    state.rotation = 0;
    state.isPlaying = false;
    if (audio) {
      audio.pause();
      audio.removeAttribute("src");
      audio.load();
    }
  }

  function seekTo(nextTime) {
    const audio = audioRef.value;
    if (!audio) {
      return;
    }

    audio.currentTime = Number(nextTime || 0);
    state.currentTime = audio.currentTime;
  }

  function handleLoadedMetadata() {
    const audio = audioRef.value;
    if (!audio) {
      return;
    }

    state.duration = Number.isFinite(audio.duration) ? audio.duration : 0;
    state.currentTime = Number.isFinite(audio.currentTime) ? audio.currentTime : 0;
  }

  function handleTimeUpdate() {
    const audio = audioRef.value;
    if (!audio) {
      return;
    }

    state.currentTime = Number.isFinite(audio.currentTime) ? audio.currentTime : 0;
  }

  function handlePlay() {
    state.isPlaying = true;
    startRotationLoop();
  }

  function handlePause() {
    state.isPlaying = false;
    stopRotationLoop();
  }

  function handleEnded() {
    state.isPlaying = false;
    stopRotationLoop();
  }

  function startRotationLoop() {
    stopRotationLoop();

    const tick = (timestamp) => {
      if (!state.isPlaying) {
        return;
      }

      if (!state.lastFrameTime) {
        state.lastFrameTime = timestamp;
      }

      const delta = timestamp - state.lastFrameTime;
      state.lastFrameTime = timestamp;
      state.rotation = (state.rotation + delta * 0.02) % 360;
      state.rafId = window.requestAnimationFrame(tick);
    };

    state.lastFrameTime = 0;
    state.rafId = window.requestAnimationFrame(tick);
  }

  function stopRotationLoop() {
    if (state.rafId) {
      window.cancelAnimationFrame(state.rafId);
      state.rafId = 0;
    }
    state.lastFrameTime = 0;
  }

  onBeforeUnmount(stopRotationLoop);

  return {
    audioRef,
    player: state,
    currentTrack: computed(() => state.currentTrack),
    discStyle,
    discLabelStyle,
    stageAmbientStyle,
    playTrack,
    togglePlay,
    syncTrack,
    clearTrack,
    seekTo,
    handleLoadedMetadata,
    handleTimeUpdate,
    handlePlay,
    handlePause,
    handleEnded,
  };
}
