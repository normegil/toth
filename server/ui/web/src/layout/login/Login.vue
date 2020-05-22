<template>
    <div class="login__container" @keyup.enter="attemptLogin">
        <InputField class="login__field--top" label="Mail" v-model="login" />
        <InputField class="login__field--middle" label="Password" v-model="password" :password="true" />
        <p v-if="showLoginError || showServerError || showUIError" class="login__field--middle login__error-text">
            <span v-if="showLoginError">Wrong mail or password</span>
            <span v-if="showServerError">Server error</span>
            <span v-if="showUIError">UI error</span>
        </p>
        <div class="login__field--bottom login__buttons-container">
            <Button class="login__button login__button--first" title="Register"/>
            <Button class="login__button login__button--last" title="Sign in" @click="attemptLogin" />
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import InputField from "../../components/InputField.vue";
import Component from "vue-class-component";
import Button from "../../components/Button.vue";
import {AxiosError} from "axios";

@Component({
    components: {Button, InputField}
})
export default class Login extends Vue {
    login = ""
    password = ""

    showLoginError = false;
    showServerError = false;
    showUIError = false;

    attemptLogin(): void {
        this.$store.dispatch("auth/signIn", {
            username: this.login,
            password: this.password
        })
            .then(() => {
                this.showLoginError = false;
                this.showServerError = false;
                this.showUIError = false;
            })
            .catch((err: AxiosError): void => {
                if (undefined !== err.response) {
                    if (err.response.status === 401) {
                        this.showLoginError = true;
                        this.showServerError = false;
                        this.showUIError = false;
                    } else if (err.response.status === 500) {
                        this.showLoginError = false;
                        this.showServerError = true;
                        this.showUIError = false;
                    }
                } else {
                    this.showUIError = true;
                    this.showLoginError = false;
                    this.showServerError = false;
                }
            });
    }
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

    &__error-text {
        border: 1px solid $color-red;
        background-color: $color-red-light;
        padding: 0.3rem 1rem;
    }
}
</style>