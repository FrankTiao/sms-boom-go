<template>
    <n-config-provider :theme="theme">
        <n-split
            direction="vertical"
            style="height: 100vh"
            :max="0.75"
            :min="0.25"
        >
            <template #1>
                <n-card>
                    <n-form ref="formRef"
                            :model="formValue"
                            :rules="formRules"
                            :show-require-mark="false"
                            label-placement="top"
                            label-width="auto"
                            size="medium"
                    >
                        <n-form-item label="手机号" path="phone">
                            <n-input
                                v-model:value="formValue.phone"
                                type="textarea"
                                placeholder="多个手机号时每行一个"
                            />
                        </n-form-item>
                        <n-form-item label="轰炸轮数" path="rounds">
                            <n-input-number
                                v-model:value="formValue.rounds"
                                placeholder="对每个手机号轰炸几轮，默认1轮"
                                :min="1"
                            />
                        </n-form-item>
                        <n-form-item label="轰炸间隔" path="interval">
                            <n-input-number
                                v-model:value="formValue.interval"
                                placeholder="每轮轰炸结束后休息几秒，默认60秒"
                                :min="0"
                            />
                        </n-form-item>
                        <n-form-item label="协程数" path="coroutineCount">
                            <n-slider
                                v-model:value="formValue.coroutineCount"
                                :min="1"
                                :max="512"
                                :step="8"
                            />
                        </n-form-item>
                        <div style="display: flex; justify-content: flex-end">
                            <n-button type="primary" @click="handleSubmit">开始轰炸</n-button>
                            <n-button style="margin-left: 10px" @click="handleReset">重置</n-button>
                        </div>
                    </n-form>
                </n-card>
            </template>
            <template #2>
                <n-card>
                    <n-scrollbar style="max-height: 300px">
                        我们在田野上面找猪<br>
                        想象中已找到了三只<br>
                        小鸟在白云上面追逐<br>
                        它们在树底下跳舞<br>
                        啦啦啦啦啦啦啦啦咧<br>
                        啦啦啦啦咧<br>
                        我们在想象中度过了许多年<br>
                        想象中我们是如此的疯狂<br>
                        我们在城市里面找猪<br>
                        想象中已找到了几百万只<br>
                        小鸟在公园里面唱歌<br>
                        它们独自在想象里跳舞<br>
                        啦啦啦啦啦啦啦啦咧<br>
                        啦啦啦啦咧<br>
                        我们在想象中度过了许多年<br>
                        许多年之后我们又开始想象<br>
                        啦啦啦啦啦啦啦啦咧
                    </n-scrollbar>
                </n-card>
            </template>
        </n-split>
    </n-config-provider>
</template>

<script setup>
import { ref, reactive, computed } from 'vue';
import { darkTheme, useOsTheme } from "naive-ui";

const osThemeRef = useOsTheme();
const theme = computed(() => osThemeRef.value === "dark" ? darkTheme : null)


const formRef = ref(null);
const formValue = reactive({
    phone: '',
    rounds: 1,
    interval: 60,
    coroutineCount: 64,
});

const formRules = {
    phone: {
        required: true,
        message: '手机号不能为空',
        trigger: 'blur',
    },
    rounds: {
        required: true,
        type: 'number',
        message: '轰炸轮数必须为数字',
        trigger: 'blur',
    },
    interval: {
        required: true,
        type: 'number',
        message: '轰炸间隔必须为数字',
        trigger: 'blur',
    },
    coroutineCount: {
        required: true,
        type: 'number',
        message: '协程数必须为数字',
        trigger: 'blur',
    },
};

const handleSubmit = () => {
    formRef.value?.validate((errors) => {
        if (!errors) {
            console.log('表单验证通过，开始轰炸！');
            // 在这里处理表单提交逻辑
            console.log(formValue);
        } else {
            console.log('表单验证失败，请检查输入！');
        }
    });
};

const handleReset = () => {
    formRef.value?.resetValidation();
    formValue.phone = '';
    formValue.rounds = 1;
    formValue.interval = 60;
    formValue.coroutineCount = 64;
};
</script>

<style scoped>
</style>
