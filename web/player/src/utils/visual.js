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

export function trackHue(track) {
  return colorSeed(track?.file_name || track?.music_name || "vmusic") % 360;
}

export function gradientStyle(track) {
  const hue = trackHue(track);
  return {
    background: `linear-gradient(135deg, hsl(${hue} 84% 72%), hsl(${(hue + 42) % 360} 74% 48%))`,
  };
}

export function ambientStyle(track) {
  const hue = trackHue(track);
  return {
    filter: `hue-rotate(${hue}deg) saturate(1.18)`,
  };
}

export function accentVars(track) {
  const hue = trackHue(track);
  return {
    "--accent-hue": `${hue}`,
    "--accent-hue-secondary": `${(hue + 36) % 360}`,
    "--accent-hue-third": `${(hue + 180) % 360}`,
  };
}

export function stageFacts(track, options = {}) {
  if (!track) {
    return [
      { label: "状态", value: "等待播放" },
      { label: "来源", value: "本地音乐库" },
      { label: "模式", value: "Glass SFC" },
    ];
  }

  return [
    { label: "来源", value: track.from || "local" },
    { label: "格式", value: String(track.format || "audio").toUpperCase() },
    { label: "收藏", value: options.liked ? "已喜欢" : "未收藏" },
  ];
}

export function stageMeta(track, options = {}) {
  if (!track) {
    return [
      "上传音频文件后，播放器会自动点亮。",
      "登录后可将歌曲加入你的 LikeList。",
      "暂停时唱片保持当前旋转角度，不会复位。",
    ];
  }

  return [
    `文件名：${track.file_name || "未知文件"}`,
    `当前状态：${options.playing ? "正在播放" : "暂停 / 待命"}`,
    `收藏状态：${options.liked ? "已加入我喜欢" : "尚未加入我喜欢"}`,
    `音频地址：${track.audio_url || "未生成"}`,
  ];
}
