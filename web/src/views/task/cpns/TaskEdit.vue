<template>
  <el-drawer title="任务添加/编辑" :visible.sync="taskEditFormShow"
             direction="rtl"
             @opened="()=>{this.$refs.mInput.focus()}"
             :before-close="handleClose"
             size="30%">
    <el-form ref="taskForm" :model="taskForm" class="drawerForm" :rules="taskFormRules"
             label-width="100px">
      <el-form-item label="任务名称" prop="task_name">
        <el-input v-model="taskForm.task_name" ref="mInput"></el-input>
      </el-form-item>
      <el-form-item label="任务说明" prop="description">
        <el-input type="textarea" :autosize="{ minRows: 4, maxRows: 16}" v-model="taskForm.description"></el-input>
      </el-form-item>
      <el-form-item label="所属分组" prop="group_id">
        <el-select v-model="taskForm.group_id">
          <el-option v-for="item in groupOptions" :key="item.group_id" :label="item.group_name" :value="item.group_id"/>
        </el-select>
      </el-form-item>
      <el-form-item label="cron表达式" prop="cron_spec">
        <el-input v-model="taskForm.cron_spec"></el-input>
      </el-form-item>
      <el-form-item label="命令脚本" prop="command">
        <el-input type="textarea" :autosize="{ minRows: 4, maxRows: 16}" v-model="taskForm.command"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button @click="handleClose">取 消</el-button>
        <el-button type="primary" @click="saveTask">确 定</el-button>
      </el-form-item>
    </el-form>
  </el-drawer>
</template>

<script>
import {saveTask} from "@/api/task";
import {getGroupList} from "@/api/group";

export default {
  name: "TaskEdit",
  data() {
    return {
      taskEditFormShow: false,
      taskForm: {
        task_id: 0,
        group_id: '',
        task_name: '',
        description: '',
        cron_spec: '',
        command: ''
      },
      groupOptions: [],
      taskFormRules: {
        task_name: [
          {required: true, message: '任务名称必填', trigger: 'blur'},
          {min: 3, max: 200, message: '长度在3到200个字符', trigger: 'blur'}
        ],
        description: [
          {required: true, message: '任务说明必填', trigger: 'blur'}
        ],
        cron_spec: [
          {required: true, message: 'cron表达式必填', trigger: 'blur'}
        ],
        command: [
          {required: true, message: '命令脚本必填', trigger: 'blur'}
        ]
      }
    }
  },
  methods: {
    setEditVal(row) {
      console.log(this.taskForm)
      for (let k in this.taskForm) {
        this.$set(this.taskForm, k, row[k] ? row[k] : '')
      }
      this.getGroups()
      this.taskEditFormShow = true
    },
    async getGroups() {
      await getGroupList({}).then((res) => {
        this.groupOptions = res.data.list
      }).catch(() => {
      })
    },
    handleClose() {
      for (let k in this.taskForm) {
        this.$set(this.taskForm, k, '')
      }
      this.$refs.taskForm.clearValidate()
      this.taskEditFormShow = false
    },
    saveTask() {
      this.$refs.taskForm.validate(async (valid) => {
        if (valid) {
          await saveTask(this.taskForm).then((res) => {
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
