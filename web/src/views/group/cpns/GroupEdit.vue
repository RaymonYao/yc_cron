<template>
  <el-drawer title="分组添加/编辑" :visible.sync="groupEditFormShow"
             direction="rtl"
             @opened="()=>{this.$refs.mInput.focus()}"
             :before-close="handleClose"
             size="30%">
    <el-form ref="groupForm" :model="groupForm" class="drawerForm" :rules="groupFormRules" label-width="70px">
      <el-form-item label="分组名称" prop="group_name">
        <el-input v-model="groupForm.group_name" ref="mInput"></el-input>
      </el-form-item>
      <el-form-item label="分组说明" prop="description">
        <el-input
            type="textarea"
            :autosize="{ minRows: 8, maxRows: 16}"
            v-model="groupForm.description">
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-button @click="handleClose">取 消</el-button>
        <el-button type="primary" @click="saveGroup">确 定</el-button>
      </el-form-item>
    </el-form>
  </el-drawer>
</template>

<script>
import {saveGroup} from "@/api/group";

export default {
  name: "GroupEdit",
  data() {
    return {
      groupEditFormShow: false,
      groupForm: {
        group_id: 0,
        group_name: '',
        description: ''
      },
      groupFormRules: {
        group_name: [
          {required: true, message: '分组名称必填', trigger: 'blur'},
          {min: 3, max: 200, message: '长度在3到200个字符', trigger: 'blur'}
        ],
        description: [
          {required: true, message: '分组说明必填', trigger: 'blur'}
        ]
      }
    }
  },
  methods: {
    setEditVal(row) {
      for (let k in this.groupForm) {
        this.$set(this.groupForm, k, row[k] ? row[k] : '')
      }
      this.groupEditFormShow = true
    },
    handleClose() {
      for (let k in this.groupForm) {
        this.$set(this.groupForm, k, '')
      }
      this.$refs.groupForm.clearValidate()
      this.groupEditFormShow = false
    },
    saveGroup() {
      this.$refs.groupForm.validate(async (valid) => {
        if (valid) {
          await saveGroup(this.groupForm).then((res) => {
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
