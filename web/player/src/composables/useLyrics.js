import { computed, ref, watch } from "vue";

function buildLyrics(track) {
  if (!track) {
    return [
      { time: 0, text: "上传一首歌，让舞台先亮起来。" },
      { time: 8, text: "这里会展示动态歌词与当前播放节奏。" },
      { time: 16, text: "后续可直接接入真实 LRC 数据。" },
    ];
  }

  const title = track.music_name || "未命名曲目";
  const author = track.author || "未知艺术家";
  return [
    { time: 0, text: `${title} 正在进入唱片舞台` },
    { time: 8, text: `来自 ${author} 的本地音乐开始播放` },
    { time: 16, text: "唱片旋转、玻璃高光和动态背景同时启动" },
    { time: 24, text: "登录后可以把这首歌加入你的我喜欢列表" },
    { time: 32, text: "当前歌词区为动态占位实现，可后续接入真实歌词" },
    { time: 40, text: `文件 ${track.file_name || "未知文件"} 正在持续播放中` },
  ];
}

export function useLyrics(currentTrack, currentTime) {
  const lines = ref(buildLyrics(currentTrack.value));

  watch(currentTrack, (track) => {
    lines.value = buildLyrics(track);
  }, { immediate: true });

  const activeIndex = computed(() => {
    const time = Number(currentTime.value || 0);
    let index = 0;
    lines.value.forEach((line, lineIndex) => {
      if (time >= line.time) {
        index = lineIndex;
      }
    });
    return index;
  });

  return {
    lines,
    activeIndex,
  };
}
