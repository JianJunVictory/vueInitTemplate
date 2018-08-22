<template>
    <div>
        <el-form :model="loginFrom" status-icon :rules="rules2" ref="loginFrom" label-width="100px" class="demo-ruleForm">
        <el-form-item :label="$t('login.email')" prop="email">
            <el-input type="email" v-model="loginFrom.email" auto-complete="off" :placeholder="$t('login.email')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('login.password')" prop="password">
            <el-input type="password" v-model="loginFrom.password" auto-complete="off" :placeholder="$t('login.password')"></el-input>
        </el-form-item>
        <el-form-item>
            <el-button type="primary" @click="submitForm('loginFrom')">{{$t('login.submit')}}</el-button>
            <el-button @click="resetForm('loginFrom')">{{$t('login.reset')}}</el-button>
        </el-form-item>
        </el-form>
    </div>
</template>

<script>
import regExpUtil from '../utils/regExpUtil'
import { mapActions } from 'vuex'
export default {
  name: 'Login',
  data () {
    var validateEmail = (rule, value, callback) => {
      if (value === '') {
        if (this.$i18n.locale === 'zh') {
          callback(new Error('请输入手机号'))
        } else {
          callback(new Error('please input phone'))
        }
      } else if (!regExpUtil.isEmail(value)) {
        callback(new Error('请输入正确手机号'))
      } else {
        callback()
      }
    }
    var validatePassword = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入密码'))
      } else if (!regExpUtil.checkPassword(value)) {
        callback(new Error('请输入至少由数字,字母组成的8到16位密码'))
      } else {
        callback()
      }
    }
    return {
      title: this.$store.state.info,
      loginFrom: {
        email: 'jianjun@dappworks.cn',
        password: '123456778'
      },
      rules2: {
        email: [
          { validator: validateEmail, trigger: 'blur' }
        ],
        password: [
          { validator: validatePassword, trigger: 'blur' }
        ]
      }
    }
  },
  methods: {
    ...mapActions(['doLogin']),
    submitForm (formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          let email = this.loginFrom.email
          let password = this.loginFrom.password
          console.log(email, password)
          // let that = this
          this.$store.dispatch('doLogin', {'email': email, 'password': password}).then(response => {
            if (response.status) {
              console.log(response)
            }
          })
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    resetForm (formName) {
      this.$refs[formName].resetFields()
    }
  }
}
</script>

<style lang="">
</style>
