import { computed, defineComponent } from "https://unpkg.com/vue@3/dist/vue.esm-browser.prod.js";
import { gradientStyle, makeInitials } from "../utils/visual.js";

export const QueuePanel = defineComponent({
  name: "QueuePanel",
  props: {
    filter: { type: String, default: "all" },
    tracks: { type: Array, default: () => [] },
    currentTrackFileName: { type: String, default: "" },
    isPlaying: { type: Boolean, default: false },
    likedFiles: { type: Array, default: () => [] },
  },
  emits: ["update:filter", "play", "toggle-like"],
  setup(props) {
    const likedSet = computed(() => new Set(props.likedFiles || []));
    const emptyText = computed(() => {
      if (props.filter === "liked") {
        return {
          title: "还没有喜欢的音乐",
          description: "登录后点击心形按钮，把喜欢的歌曲收进你的 LikeList。",
        };
      }
      return {
        title: "当前音乐库为空",
        description: "先上传一首歌，播放器舞台就会自动点亮。",
      };
    });

    function isLiked(track) {
      return likedSet.value.has(track.file_name);
    }

    function coverStyle(track) {
      return gradientStyle(track);
    }

    function initials(track) {
      return makeInitials(track.music_name || track.file_name || "VM");
    }

    return {
      emptyText,
      isLiked,
      coverStyle,
      initials,
    };
  },
  template: `
    <section class="queue-panel glass-panel">
      <div class="panel-head panel-head--queue">
        <div>
          <p class="eyebrow">歌单与收藏</p>
          <h3>本地音乐队列</h3>
        </div>
        <div class="segmented segmented--soft">
          <button type="button" :class="{ active: filter === 'all' }" @click="$emit('update:filter', 'all')">全部</button>
          <button type="button" :class="{ active: filter === 'liked' }" @click="$emit('update:filter', 'liked')">我喜欢</button>
        </div>
      </div>

      <div v-if="tracks.length" class="queue-list">
        <article v-for="track in tracks" :key="track.file_name" class="queue-item" :class="{ active: currentTrackFileName === track.file_name }">
          <button class="queue-item__main" type="button" @click="$emit('play', track)">
            <span class="queue-item__cover" :style="coverStyle(track)">{{ initials(track) }}</span>
            <span class="queue-item__copy">
              <strong>{{ track.music_name }}</strong>
              <span>{{ track.author }}</span>
            </span>
          </button>

          <div class="queue-item__meta">
            <span class="queue-item__tag">{{ currentTrackFileName === track.file_name && isPlaying ? '播放中' : String(track.format || 'audio').toUpperCase() }}</span>
            <button class="icon-button" :class="{ active: isLiked(track) }" type="button" @click="$emit('toggle-like', track)">{{ isLiked(track) ? '♥' : '♡' }}</button>
          </div>
        </article>
      </div>

      <div v-else class="empty-state empty-state--soft">
        <strong>{{ emptyText.title }}</strong>
        <span>{{ emptyText.description }}</span>
      </div>
    </section>
  `,
});
