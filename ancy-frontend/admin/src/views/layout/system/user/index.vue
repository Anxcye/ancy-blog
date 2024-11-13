<template>
  <div class="app-container">
    <el-form ref="queryForm" :model="queryParams" class="query-form" size="small">
      <el-input
        v-model="queryParams.userName"
        placeholder="用户名"
        clearable
        @keyup.enter="getUserList()"
      />
      <el-select v-model="queryParams.status" placeholder="状态" clearable>
        <el-option key="0" label="正常" value="0" />
        <el-option key="1" label="停用" value="1" />
      </el-select>
      <el-button type="primary" :icon="Search" @click="getUserList()">搜索</el-button>
    </el-form>

    <el-row :gutter="10">
      <el-col :span="1.5">
        <el-button type="primary" plain :icon="Plus" @click="handleAdd">新增</el-button>
      </el-col>
    </el-row>

    <el-table :data="userList" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="50" align="center" />
      <el-table-column prop="userName" label="用户名" align="center" />
      <el-table-column prop="nickName" label="昵称" align="center" />
      <el-table-column prop="roleIds" label="角色" align="center">
        <template v-slot="scope">
          <span v-for="roleId in scope.row.roleIds" :key="roleId">
            <el-tag>{{ roleList.find((role) => role.id === roleId)?.roleName ?? roleId }}</el-tag>
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" align="center">
        <template v-slot="scope">
          <el-switch
            v-model="scope.row.status"
            active-value="0"
            inactive-value="1"
            @change="handleStatusChange(scope.row)"
            :loading="statusLoading"
          />
        </template>
      </el-table-column>
      <el-table-column label="时间" align="center" width="250" prop="createTime">
        <template v-slot="scope">
          <div>创建时间: {{ scope.row.createTime }}</div>
          <div v-if="scope.row.updateTime">更新时间: {{ scope.row.updateTime }}</div>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" fixed="right" width="170">
        <template v-slot="scope">
          <el-button type="text" :icon="Edit" @click="handleUpdate(scope.row)" size="small">
            修改
          </el-button>
          <el-button
            size="small"
            type="text"
            :icon="Delete"
            @click="handleDelete(scope.row)"
            :disabled="scope.row.id === 1"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:page-size="queryParams.pageSize"
      layout="sizes, prev, pager, next, jumper, ->, total"
      :total="total"
      :page-sizes="[5, 10, 20, 30, 40]"
      v-model:current-page="queryParams.pageNum"
      @current-change="getUserList"
      @size-change="getUserList()"
    />

    <!-- 添加或修改参数配置对话框 -->
    <el-dialog :title="title" v-model="open" class="dialog-form">
      <el-form ref="form" :model="user" :rules="rules" label-width="auto" label-position="top">
        <el-form-item label="用户昵称" prop="nickName">
          <el-input v-model="user.nickName" placeholder="用户昵称" maxlength="30" />
        </el-form-item>
        <!-- <el-form-item label="邮箱" prop="email">
          <el-input v-model="user.email" placeholder="邮箱" maxlength="50" />
        </el-form-item> -->
        <el-form-item label="用户名" prop="userName">
          <el-input v-model="user.userName" placeholder="用户名" maxlength="30" />
        </el-form-item>
        <el-form-item label="用户密码" prop="password">
          <el-input
            v-model="user.password"
            placeholder="请输入用户密码"
            type="password"
            maxlength="20"
            show-password
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="user.status">
            <el-radio key="0" label="0">正常</el-radio>
            <el-radio key="1" label="1">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="user.roleIds" multiple placeholder="请选择">
            <el-option
              v-for="item in roleList"
              :key="item.id"
              :label="item.roleName"
              :value="item.id"
              :disabled="item.status == '1'"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template v-slot:footer>
        <div class="dialog-footer">
          <el-button @click="handleCancel">取 消</el-button>
          <el-button type="primary" @click="handleSubmit">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reqUserAdd, reqUserDelete, reqUserPage, reqUserUpdate } from '@/api/system/user'
import type { UserAddParams, UserListData, UserPageParams } from '@/api/system/user/type'
import { ref, onMounted } from 'vue'
import { Plus, Edit, Delete, Search } from '@element-plus/icons-vue'
import type { FormRules } from 'element-plus'
import { reqRoleList } from '@/api/system/role'
import type { RoleListData } from '@/api/system/role/type'
import { toggleStatus } from '@/utils/toggleStatus'

const queryParams = ref<UserPageParams>({
  pageNum: 1,
  pageSize: 10,
})
const userList = ref<UserListData[]>([])
const total = ref<number>(0)
const loading = ref<boolean>(false)
const open = ref<boolean>(false)
const title = ref<string>('')
const user = ref<UserAddParams>({})
const rules = ref<FormRules>({})
const statusLoading = ref<boolean>(false)

const roleList = ref<RoleListData[]>([])

const getUserList = async (page: number = 1) => {
  queryParams.value.pageNum = page
  loading.value = true
  const res = await reqUserPage(queryParams.value)
  userList.value = res.data.rows
  total.value = res.data.total
  loading.value = false
}

const getRoleList = async () => {
  const res = await reqRoleList()
  roleList.value = res.data
}

const handleAdd = () => {
  user.value = {}
  open.value = true
  title.value = '新增'
}

const handleStatusChange = async (row: UserListData) => {
  statusLoading.value = true
  try {
    await reqUserUpdate(row.id, { status: row.status })
  } catch {
    toggleStatus(row)
  } finally {
    statusLoading.value = false
  }
}

const handleUpdate = (row: UserListData) => {
  user.value = row
  open.value = true
  title.value = '修改' + row.nickName
}

const handleDelete = (row: UserListData) => {
  ElMessageBox.confirm('删除?', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
  }).then(async () => {
    await reqUserDelete(row.id)
    await getUserList()
  })
}

const handleSubmit = async () => {
  if (user.value.id) {
    await reqUserUpdate(user.value.id, user.value)
  } else {
    await reqUserAdd(user.value)
    ElMessage.success('新增成功')
  }
  open.value = false
  getUserList()
}

const handleCancel = () => {
  user.value = {}
  open.value = false
}

onMounted(() => {
  getUserList()
  getRoleList()
})
</script>

<style scoped lang="scss"></style>
