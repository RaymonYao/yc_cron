<template>
  <el-drawer title="节点添加/编辑" :visible.sync="workerEditFormShow"
             direction="rtl"
             @opened="()=>{this.$refs.mInput.focus()}"
             :before-close="handleClose"
             size="30%">
    <el-form ref="workerForm" :model="workerForm" class="drawerForm" :rules="workerFormRules" label-width="80px">
      <el-form-item label="节点IP" prop="worker_ip">
        <el-input v-model="workerForm.worker_ip" ref="mInput"></el-input>
      </el-form-item>
      <el-form-item label="节点说明" prop="description">
        <el-input
            type="textarea"
            :autosize="{ minRows: 8, maxRows: 16}"
            v-model="workerForm.description">
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-button @click="handleClose">取 消</el-button>
        <el-button type="primary" @click="saveWorker">确 定</el-button>
      </el-form-item>
    </el-form>
  </el-drawer>
</template>

<script>
import {saveWorker} from "@/api/worker";

export default {
  name: "WorkerEdit",
  data() {
    return {
      workerEditFormShow: false,
      workerForm: {
        worker_id: 0,
        worker_ip: '',
        description: ''
      },
      workerFormRules: {
        worker_ip: [
          {required: true, message: '节点IP必填', trigger: 'blur'}
        ],
        description: [
          {required: true, message: '分组说明必填', trigger: 'blur'}
        ]
      }
    }
  },
  methods: {
    setEditVal(row) {
      for (let k in this.workerForm) {
        this.$set(this.workerForm, k, row[k] ? row[k] : '')
      }
      this.workerEditFormShow = true
    },
    handleClose() {
      for (let k in this.workerForm) {
        this.$set(this.workerForm, k, '')
      }
      this.$refs.workerForm.clearValidate()
      this.workerEditFormShow = false
    },
    saveWorker() {
      this.$refs.workerForm.validate(async (valid) => {
        if (valid) {
          await saveWorker(this.workerForm).then((res) => {
            this.$message({
              type: 'success',
              message: res.msg,
              showClose: true
            });
            this.$parent.getTableData()
            this.handleClose()
          }).catch(() => {
          })
        }
      })
    }
  }
}
</script>
