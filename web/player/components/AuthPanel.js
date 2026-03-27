import {
  computed,
  defineComponent,
  reactive,
  watch,
} from "https://unpkg.com/vue@3/dist/vue.esm-browser.prod.js";

export const AuthPanel = defineComponent({
  name: "AuthPanel",
  props: {
    session: { type: Object, required: true },
    mode: { type: String, default: "login" },
    loading: { type: Boolean, default: false },
    likedCount: { type: Number, default: 0 },
    totalTracks: { type: Number, default: 0 },
    userInitials: { type: String, default: "VM" },
  },
  emits: ["update:mode", "login", "register", "logout"],
  setup(props, { emit }) {
    const loginForm = reactive({
      username: "",
      password: "",
    });

    const registerForm = reactive({
      username: "",
      password: "",
      mail: "",
      telephone: "",
    });

    const isLogin = computed(() => props.mode === "login");

    function switchMode(nextMode) {
      emit("update:mode", nextMode);
    }

    function submitLogin() {
      emit("login", { ...loginForm });
    }

    function submitRegister() {
      emit("register", { ...registerForm });
    }

    function resetForms() {
      loginForm.username = "";
      loginForm.password = "";
      registerForm.username = "";
      registerForm.password = "";
      registerForm.mail = "";
      registerForm.telephone = "";
    }

    watch(() => props.session.user && props.session.user.userid, (nextUserId) => {
      if (nextUserId) {
        resetForms();
      }
    });

    return {
      loginForm,
      registerForm,
      isLogin,
      switchMode,
      submitLogin,
      submitRegister,
    };
  },
  template: `
    <section class="auth-panel glass-panel">
      <div class="panel-head">
        <div>
          <p class="eyebrow">账号中心</p>
          <h3>{{ session.user ? '我的音乐身份' : '登录 / 注册' }}</h3>
        </div>
        <div v-if="!session.user" class="segmented segmented--soft">
          <button type="button" :class="{ active: mode === 'login' }" @click="switchMode('login')">登录</button>
          <button type="button" :class="{ active: mode === 'register' }" @click="switchMode('register')">注册</button>
        </div>
      </div>

      <template v-if="session.user">
        <div class="profile-card">
          <div class="profile-card__avatar">{{ userInitials }}</div>
          <div class="profile-card__meta">
            <strong>{{ session.user.username }}</strong>
            <span>{{ session.user.mail || '未设置邮箱' }}</span>
            <span>{{ session.user.telephone || '未设置手机号' }}</span>
          </div>
        </div>

        <div class="profile-stats">
          <div>
            <span class="chip-label">我喜欢</span>
            <strong>{{ likedCount }}</strong>
          </div>
          <div>
            <span class="chip-label">全部音乐</span>
            <strong>{{ totalTracks }}</strong>
          </div>
        </div>

        <button class="ghost-button" type="button" @click="$emit('logout')">退出当前账号</button>
      </template>

      <form v-else-if="isLogin" class="auth-form" @submit.prevent="submitLogin">
        <label>
          <span>用户名</span>
          <input v-model.trim="loginForm.username" type="text" maxlength="32" placeholder="输入用户名">
        </label>
        <label>
          <span>密码</span>
          <input v-model="loginForm.password" type="password" maxlength="128" placeholder="输入密码">
        </label>
        <button class="primary-button" type="submit" :disabled="loading">{{ loading ? '登录中...' : '登录' }}</button>
      </form>

      <form v-else class="auth-form" @submit.prevent="submitRegister">
        <label>
          <span>用户名</span>
          <input v-model.trim="registerForm.username" type="text" maxlength="32" placeholder="创建用户名">
        </label>
        <label>
          <span>密码</span>
          <input v-model="registerForm.password" type="password" maxlength="128" placeholder="建议至少 6 位">
        </label>
        <label>
          <span>邮箱</span>
          <input v-model.trim="registerForm.mail" type="email" maxlength="64" placeholder="可选">
        </label>
        <label>
          <span>手机号</span>
          <input v-model.trim="registerForm.telephone" type="text" maxlength="20" placeholder="可选">
        </label>
        <button class="primary-button" type="submit" :disabled="loading">{{ loading ? '注册中...' : '注册并登录' }}</button>
      </form>
    </section>
  `,
});
