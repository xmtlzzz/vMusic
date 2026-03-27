import {
  computed,
  defineComponent,
  onMounted,
  reactive,
  ref,
} from "https://unpkg.com/vue@3/dist/vue.esm-browser.prod.js";

import { AuthPanel } from "./components/AuthPanel.js";
import { HeaderBar } from "./components/HeaderBar.js";
import { PlayerStage } from "./components/PlayerStage.js";
import { QueuePanel } from "./components/QueuePanel.js";
import { UploadPanel } from "./components/UploadPanel.js";
import {
  fetchCurrentUser,
  fetchTracks,
  loginUser,
  registerUser,
  removeToken,
  toggleLike,
  uploadTracks,
} from "./services/api.js";
import { usePlayer } from "./composables/usePlayer.js";
import { makeInitials } from "./utils/visual.js";

export const AppRoot = defineComponent({
  name: "AppRoot",
  components: {
    AuthPanel,
    HeaderBar,
    PlayerStage,
    QueuePanel,
    UploadPanel,
  },
  setup() {
    const tracks = ref([]);
    const libraryPath = ref("storage/music");

    const session = reactive({
      token: localStorage.getItem("vmusic_access_token") || "",
      user: null,
    });

    const ui = reactive({
      authMode: "login",
      filter: "all",
      notice: "",
      error: "",
      authLoading: false,
      uploading: false,
    });

    const {
      audioRef,
      player,
      currentTrack,
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
    } = usePlayer({
      onAutoplayBlocked: () => {
        setError("浏览器阻止了自动播放，请再次点击播放。");
      },
    });

    const likedFiles = computed(() => {
      const list = session.user && Array.isArray(session.user.like_list) ? session.user.like_list : [];
      return list.map((item) => item && item.file_name).filter(Boolean);
    });

    const likedSet = computed(() => new Set(likedFiles.value));
    const likedCount = computed(() => likedFiles.value.length);
    const displayedTracks = computed(() => {
      if (ui.filter === "liked") {
        return tracks.value.filter((track) => likedSet.value.has(track.file_name));
      }
      return tracks.value;
    });
    const activeQueue = computed(() => displayedTracks.value.length ? displayedTracks.value : tracks.value);
    const currentTrackLiked = computed(() => {
      return !!(currentTrack.value && likedSet.value.has(currentTrack.value.file_name));
    });
    const userInitials = computed(() => makeInitials(session.user && session.user.username ? session.user.username : "VM"));
    const libraryPathText = computed(() => `音乐库目录：${libraryPath.value || "storage/music"}`);

    onMounted(async () => {
      await loadTracks();
      await hydrateSession();
    });

    async function loadTracks() {
      clearError();
      try {
        const payload = await fetchTracks();
        tracks.value = Array.isArray(payload.music_list) ? payload.music_list : [];
        libraryPath.value = payload.library_path || "storage/music";

        if (!currentTrack.value) {
          if (tracks.value.length) {
            await playTrack(tracks.value[0], { autoplay: false });
          }
          return;
        }

        const matchedTrack = tracks.value.find((track) => track.file_name === currentTrack.value.file_name);
        if (matchedTrack) {
          syncTrack(matchedTrack);
          return;
        }

        if (tracks.value.length) {
          await playTrack(tracks.value[0], { autoplay: false });
          return;
        }

        clearTrack();
      } catch (error) {
        setError(error.message);
      }
    }

    async function hydrateSession() {
      if (!session.token) {
        return;
      }

      try {
        session.user = await fetchCurrentUser(session.token);
      } catch (_) {
        clearSession();
        setError("登录态已失效，请重新登录。");
      }
    }

    async function submitLogin(credentials) {
      if (!credentials.username || !credentials.password) {
        setError("请输入用户名和密码。");
        return;
      }

      ui.authLoading = true;
      clearError();
      try {
        const payload = await loginUser(credentials);
        applyAuth(payload);
        ui.notice = payload.message || "登录成功。";
      } catch (error) {
        setError(error.message);
      } finally {
        ui.authLoading = false;
      }
    }

    async function submitRegister(payload) {
      if (!payload.username || !payload.password) {
        setError("用户名和密码不能为空。");
        return;
      }

      ui.authLoading = true;
      clearError();
      try {
        await registerUser(payload);
        const authPayload = await loginUser({ username: payload.username, password: payload.password });
        applyAuth(authPayload);
        ui.notice = "注册完成并已自动登录。";
      } catch (error) {
        setError(error.message);
      } finally {
        ui.authLoading = false;
      }
    }

    async function logout() {
      const username = session.user ? session.user.username : "";
      const accessToken = session.token;
      clearSession();

      if (!username || !accessToken) {
        return;
      }

      try {
        await removeToken(username, accessToken);
      } catch (_) {
        // Ignore backend cleanup errors after local session is cleared.
      }
    }

    async function uploadSelectedFiles(files, resetSelection) {
      if (!files.length) {
        setError("请先选择要上传的音频文件。");
        return;
      }

      ui.uploading = true;
      clearError();
      try {
        const payload = await uploadTracks(files);
        ui.notice = `上传完成：${payload.saved}/${payload.total} 首。`;
        await loadTracks();
        if (typeof resetSelection === "function") {
          resetSelection();
        }

        if (Array.isArray(payload.music_list) && payload.music_list.length) {
          const latestTrack = payload.music_list[payload.music_list.length - 1];
          const matchedTrack = tracks.value.find((track) => track.file_name === latestTrack.file_name);
          if (matchedTrack) {
            await playTrack(matchedTrack, { autoplay: true });
          }
        }
      } catch (error) {
        setError(error.message);
      } finally {
        ui.uploading = false;
      }
    }

    async function handlePlayTrack(track) {
      await playTrack(track, { autoplay: true });
    }

    async function handleTogglePlay() {
      const fallbackTrack = activeQueue.value.length ? activeQueue.value[0] : null;
      await togglePlay(fallbackTrack);
    }

    async function stepTrack(direction) {
      if (!activeQueue.value.length) {
        return;
      }

      const currentIndex = currentTrack.value
        ? activeQueue.value.findIndex((track) => track.file_name === currentTrack.value.file_name)
        : -1;

      const nextIndex = currentIndex === -1
        ? 0
        : (currentIndex + direction + activeQueue.value.length) % activeQueue.value.length;

      await playTrack(activeQueue.value[nextIndex], { autoplay: true });
    }

    async function toggleTrackLike(track = currentTrack.value) {
      if (!track) {
        return;
      }
      if (!session.user || !session.token) {
        ui.authMode = "login";
        setError("请先登录，再把喜欢的歌曲加入你的列表。");
        return;
      }

      clearError();
      try {
        const payload = await toggleLike(session.token, track.file_name);
        if (payload.user) {
          session.user = payload.user;
        } else if (session.user) {
          session.user.like_list = Array.isArray(payload.like_list) ? payload.like_list : [];
        }
        ui.notice = payload.liked
          ? `已加入“我喜欢”：${track.music_name}`
          : `已取消喜欢：${track.music_name}`;
      } catch (error) {
        setError(error.message);
      }
    }

    function handleSeek(nextTime) {
      seekTo(nextTime);
    }

    async function handleTrackEnded() {
      handleEnded();
      await stepTrack(1);
    }

    function applyAuth(payload) {
      session.user = payload.user || null;
      session.token = payload.token ? payload.token.access_token : "";
      persistSession();
    }

    function persistSession() {
      if (session.token) {
        localStorage.setItem("vmusic_access_token", session.token);
      } else {
        localStorage.removeItem("vmusic_access_token");
      }
    }

    function clearSession() {
      session.user = null;
      session.token = "";
      persistSession();
    }

    function setError(message) {
      ui.error = message || "";
      if (message) {
        ui.notice = "";
      }
    }

    function clearError() {
      ui.error = "";
    }

    return {
      audioRef,
      session,
      ui,
      player,
      currentTrack,
      discStyle,
      discLabelStyle,
      stageAmbientStyle,
      tracks,
      displayedTracks,
      activeQueue,
      likedFiles,
      likedCount,
      currentTrackLiked,
      userInitials,
      libraryPathText,
      loadTracks,
      submitLogin,
      submitRegister,
      logout,
      uploadSelectedFiles,
      handlePlayTrack,
      handleTogglePlay,
      stepTrack,
      toggleTrackLike,
      handleSeek,
      handleLoadedMetadata,
      handleTimeUpdate,
      handlePlay,
      handlePause,
      handleTrackEnded,
    };
  },
  template: `
    <div class="page-shell">
      <div class="page-orb page-orb--one"></div>
      <div class="page-orb page-orb--two"></div>
      <div class="page-grid"></div>

      <div class="app-shell">
        <HeaderBar
          :session="session"
          :total-tracks="tracks.length"
          :liked-count="likedCount"
          :current-track="currentTrack"
          @logout="logout"
        />

        <div v-if="ui.notice" class="banner banner--notice">{{ ui.notice }}</div>
        <div v-if="ui.error" class="banner banner--error">{{ ui.error }}</div>

        <main class="dashboard-layout">
          <section class="dashboard-main">
            <PlayerStage
              :current-track="currentTrack"
              :player="player"
              :disc-style="discStyle"
              :disc-label-style="discLabelStyle"
              :ambient-style="stageAmbientStyle"
              :current-track-liked="currentTrackLiked"
              :queue-count="activeQueue.length"
              @toggle-like="toggleTrackLike()"
              @toggle-play="handleTogglePlay"
              @previous="stepTrack(-1)"
              @next="stepTrack(1)"
              @seek="handleSeek"
            />

            <UploadPanel
              :uploading="ui.uploading"
              :library-path-text="libraryPathText"
              @upload="uploadSelectedFiles"
              @refresh="loadTracks"
            />
          </section>

          <aside class="dashboard-side">
            <AuthPanel
              :session="session"
              :mode="ui.authMode"
              :loading="ui.authLoading"
              :liked-count="likedCount"
              :total-tracks="tracks.length"
              :user-initials="userInitials"
              @update:mode="ui.authMode = $event"
              @login="submitLogin"
              @register="submitRegister"
              @logout="logout"
            />

            <QueuePanel
              :filter="ui.filter"
              :tracks="displayedTracks"
              :current-track-file-name="currentTrack ? currentTrack.file_name : ''"
              :is-playing="player.isPlaying"
              :liked-files="likedFiles"
              @update:filter="ui.filter = $event"
              @play="handlePlayTrack"
              @toggle-like="toggleTrackLike"
            />
          </aside>
        </main>

        <audio
          ref="audioRef"
          preload="metadata"
          @loadedmetadata="handleLoadedMetadata"
          @timeupdate="handleTimeUpdate"
          @play="handlePlay"
          @pause="handlePause"
          @ended="handleTrackEnded"
        ></audio>
      </div>
    </div>
  `,
});






