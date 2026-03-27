import { computed, nextTick, onBeforeUnmount, reactive, ref } from "vue";

import { ambientStyle, gradientStyle } from "../utils/visual";

export function usePlayer(options = {}) {
  const audioRef = ref(null);
  const player = reactive({
    currentTrack: null,
    isPlaying: false,
    currentTime: 0,
    duration: 0,
    rotation: 0,
    rafId: 0,
    lastFrameTime: 0,
  });

  const currentTrack = computed(() => player.currentTrack);
  const discStyle = computed(() => ({
    transform: `rotate(${player.rotation}deg)`,
  }));
  const discLabelStyle = computed(() => gradientStyle(player.currentTrack));
  const stageAmbientStyle = computed(() => ambientStyle(player.currentTrack));

  async function playTrack(track, config = {}) {
    if (!track) {
      return;
    }

    const autoplay = Boolean(config.autoplay);
    const audio = audioRef.value;
    if (!audio) {
      return;
    }

    const sameTrack = player.currentTrack && player.currentTrack.file_name === track.file_name;
    if (!sameTrack) {
      stopRotationLoop();
      player.rotation = 0;
      player.currentTrack = track;
      player.currentTime = 0;
      player.duration = 0;
      player.isPlaying = false;
      audio.src = track.audio_url;
      audio.load();
    } else {
      player.currentTrack = track;
    }

    if (autoplay) {
      try {
        await nextTick();
        await audio.play();
      } catch (error) {
        options.onAutoplayBlocked?.(error);
      }
    }
  }

  async function togglePlay(fallbackTrack) {
    const audio = audioRef.value;
    if (!audio) {
      return;
    }

    if (!player.currentTrack && fallbackTrack) {
      await playTrack(fallbackTrack, { autoplay: true });
      return;
    }

    if (!player.currentTrack) {
      return;
    }

    if (audio.paused) {
      try {
        await audio.play();
      } catch (error) {
        options.onAutoplayBlocked?.(error);
      }
      return;
    }

    audio.pause();
  }

  function syncTrack(track) {
    if (player.currentTrack && track && player.currentTrack.file_name === track.file_name) {
      player.currentTrack = track;
    }
  }

  function seekTo(nextTime) {
    const audio = audioRef.value;
    if (!audio) {
      return;
    }

    audio.currentTime = Number(nextTime || 0);
    player.currentTime = audio.currentTime;
  }

  function clearTrack() {
    stopRotationLoop();
    player.currentTrack = null;
    player.currentTime = 0;
    player.duration = 0;
    player.rotation = 0;
    player.isPlaying = false;
    const audio = audioRef.value;
    if (audio) {
      audio.pause();
      audio.removeAttribute("src");
      audio.load();
    }
  }

  function handleLoadedMetadata() {
    const audio = audioRef.value;
    if (!audio) {
      return;
    }
    player.duration = Number.isFinite(audio.duration) ? audio.duration : 0;
    player.currentTime = Number.isFinite(audio.currentTime) ? audio.currentTime : 0;
  }

  function handleTimeUpdate() {
    const audio = audioRef.value;
    if (!audio) {
      return;
    }
    player.currentTime = Number.isFinite(audio.currentTime) ? audio.currentTime : 0;
  }

  function handlePlay() {
    player.isPlaying = true;
    startRotationLoop();
  }

  function handlePause() {
    player.isPlaying = false;
    stopRotationLoop();
  }

  function handleEnded() {
    player.isPlaying = false;
    stopRotationLoop();
  }

  function startRotationLoop() {
    stopRotationLoop();
    const tick = (timestamp) => {
      if (!player.isPlaying) {
        return;
      }
      if (!player.lastFrameTime) {
        player.lastFrameTime = timestamp;
      }
      const delta = timestamp - player.lastFrameTime;
      player.lastFrameTime = timestamp;
      player.rotation = (player.rotation + delta * 0.02) % 360;
      player.rafId = window.requestAnimationFrame(tick);
    };

    player.lastFrameTime = 0;
    player.rafId = window.requestAnimationFrame(tick);
  }

  function stopRotationLoop() {
    if (player.rafId) {
      window.cancelAnimationFrame(player.rafId);
      player.rafId = 0;
    }
    player.lastFrameTime = 0;
  }

  onBeforeUnmount(stopRotationLoop);

  return {
    audioRef,
    player,
    currentTrack,
    discStyle,
    discLabelStyle,
    stageAmbientStyle,
    playTrack,
    togglePlay,
    syncTrack,
    seekTo,
    clearTrack,
    handleLoadedMetadata,
    handleTimeUpdate,
    handlePlay,
    handlePause,
    handleEnded,
  };
}
