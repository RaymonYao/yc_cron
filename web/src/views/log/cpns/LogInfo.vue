<template>
  <div class="log_info">
    <el-divider content-position="left">输出日志</el-divider>
    <div style="margin: 0 15px 10px 15px">
      <span>输出:</span>
      <div style="width: 100%;overflow: scroll;">
        <pre>{{ out_put }}</pre>
      </div>
      <span>错误信息:</span> {{ error }}
    </div>
  </div>
</template>

<script>
export default {
  name: "TaskInfo",
  props: {
    out_put: String,
    error: String
  }
}
</script>

<style scoped>
.log_info > div {
  border: 1px solid #EEE;
  margin-bottom: 5px;
  border-radius: 3px;
  font-size: 0.8em;
}

.log_info span {
  font-weight: bolder;
  color: coral;
  font-size: 0.8em;
}

.log_info pre {
  font-size: 0.8em;
}
</style>
