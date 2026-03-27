export function formatTime(value) {
  const total = Number.isFinite(value) && value > 0 ? value : 0;
  const minutes = Math.floor(total / 60);
  const seconds = Math.floor(total % 60);
  return `${String(minutes).padStart(2, "0")}:${String(seconds).padStart(2, "0")}`;
}

export function makeInitials(text) {
  const normalized = String(text || "VM").trim();
  return Array.from(normalized).slice(0, 2).join("").toUpperCase() || "VM";
}

export function colorSeed(value) {
  return Array.from(String(value || "vmusic")).reduce((sum, char) => sum + char.charCodeAt(0), 0);
}

export function gradientStyle(track) {
  const seed = colorSeed(track && track.file_name ? track.file_name : "vmusic");
  const hue = seed % 360;
  return {
    background: `linear-gradient(135deg, hsl(${hue} 82% 70%), hsl(${(hue + 36) % 360} 72% 48%))`,
  };
}

export function ambientStyle(track) {
  const seed = colorSeed(track && track.file_name ? track.file_name : "vmusic");
  const hue = seed % 360;
  return {
    filter: `hue-rotate(${hue}deg) saturate(1.15)`,
  };
}

export function buildSceneChips(track, options = {}) {
  if (!track) {
    return [
      { label: "状态", value: "等待播放" },
      { label: "来源", value: "本地音乐库" },
      { label: "模式", value: "Glass UI" },
    ];
  }

  return [
    { label: "来源", value: track.from || "local" },
    { label: "格式", value: String(track.format || "audio").toUpperCase() },
    { label: "喜欢", value: options.liked ? "已收藏" : "未收藏" },
  ];
}

export function buildStageRows(track, options = {}) {
  if (!track) {
    return [
      "上传音频文件，立即进入播放舞台。",
      "登录后可将歌曲加入你的“我喜欢”。",
      "唱片舞台保持旋转角度，暂停时不回弹。",
    ];
  }

  const rows = [
    `文件名：${track.file_name || "未知文件"}`,
    `当前状态：${options.playing ? "正在播放" : "待命"}`,
    `收藏状态：${options.liked ? "已加入 LikeList" : "还未加入 LikeList"}`,
  ];

  if (track.audio_url) {
    rows.push(`播放地址：${track.audio_url}`);
  }

  return rows;
}
