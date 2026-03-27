<script setup>
import { computed, ref } from "vue";

const props = defineProps({
  uploading: { type: Boolean, default: false },
  libraryPathText: { type: String, default: "音乐库目录：storage/music" },
});

const emit = defineEmits(["upload", "refresh"]);
const inputRef = ref(null);
const selectedNames = ref([]);

const summaryText = computed(() => {
  if (!selectedNames.value.length) {
    return "支持 mp3、wav、flac、m4a、ogg。选择文件后可直接上传。";
  }
  const preview = selectedNames.value.slice(0, 3).join("、");
  return `已选择 ${selectedNames.value.length} 个文件：${preview}${selectedNames.value.length > 3 ? ' ...' : ''}`;
});

function handleChange(event) {
  selectedNames.value = Array.from(event.target.files || []).map((file) => file.name);
}

function handleSubmit() {
  const files = Array.from(inputRef.value?.files || []);
  emit("upload", files, () => {
    selectedNames.value = [];
    if (inputRef.value) {
      inputRef.value.value = "";
    }
  });
}
</script>

<template>
  <section class="upload-panel glass-panel">
    <div class="panel-head">
      <div>
        <p class="eyebrow">上传音乐</p>
        <h3>把本地音频拖进云感舞台</h3>
      </div>
      <button class="ghost-button compact" type="button" @click="emit('refresh')">刷新歌单</button>
    </div>

    <label class="upload-dropzone upload-dropzone--glass">
      <input ref="inputRef" type="file" accept="audio/*,.mp3,.wav,.flac,.m4a,.ogg" multiple @change="handleChange">
      <span class="upload-dropzone__icon">♪</span>
      <strong>选择音频文件</strong>
      <span>{{ summaryText }}</span>
    </label>

    <div class="upload-panel__footer">
      <span class="subtle">{{ libraryPathText }}</span>
      <button class="primary-button" type="button" @click="handleSubmit" :disabled="uploading">{{ uploading ? '上传中...' : '上传到音乐库' }}</button>
    </div>
  </section>
</template>
