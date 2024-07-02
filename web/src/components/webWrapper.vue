<template>
  <div class="web-wrapper">
    <iframe :src="url" class="iframe"></iframe>
  </div>
</template>

<script>
export default {
  name: 'webWrapper',
  props: {
    url: {
      type: String,
      required: true
    }
  }
}
</script>

<style scoped>
.web-wrapper {
  width: 100%;
  height: 100vh;
  overflow: hidden;
}

.iframe {
  width: 100%;
  height: 100%;
  border: none;
}
</style>
