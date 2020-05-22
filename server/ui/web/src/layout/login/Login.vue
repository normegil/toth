<template>
    <div class="login__container">
        <InputField class="login__field--top" label="Mail" v-model="login" />
        <InputField class="login__field--middle" label="Password" v-model="password" :password="true" />
        <div class="login__field--bottom login__buttons-container">
            <Button class="login__button login__button--first" title="Register"/>
            <Button class="login__button login__button--last" title="Sign in"/>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import InputField from "../../components/InputField.vue";
import Component from "vue-class-component";
import Button from "../../components/Button.vue";

@Component({
    components: {Button, InputField}
})
export default class Login extends Vue {
    login = ""
    password = ""
}
</script>

<style lang="scss">
.login {
    &__container {
        display: flex;
        justify-content: space-between;
        flex-direction: column;
    }

    &__buttons-container {
        display: flex;
        justify-content: space-between;
    }

    &__field {
        &--top {
            margin-bottom: 0.5rem;
        }
        &--middle {
            margin-top: 0.5rem;
            margin-bottom: 0.5rem;
        }
        &--bottom {
            margin-top: 0.5rem;
        }
    }

    &__button {
        flex-grow: 1;

        &--first {
            margin-right: 0.5rem;
        }
        &--last {
            margin-left: 0.5rem;
        }
    }
}
</style>