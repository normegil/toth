<template>
    <Frame class="passwords__container">
        <h2 class="passwords__title">Mot de passe</h2>
        <div class="passwords__form">
            <InputField class="passwords__field" label="Mot de passe actuel" :password="true" v-model="currentPassword" />
            <InputField class="passwords__field" label="Nouveau mot de passe" :password="true" v-model="newPassword" />
            <InputField class="passwords__field" label="Répéter le nouveau mot de passe" :password="true" v-model="repeatedPassword" />
            <p v-if="showPasswordChangeSuccess" class="passwords__field passwords__field-message--success">
                <span>Password change successful</span>
            </p>
            <p v-if="showError" class="passwords__field passwords__field-message--error">
                <span>{{ errorMessage }}</span>
            </p>
            <Button title="Changer le mot de passe" @click="updatePassword"/>
        </div>
    </Frame>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Frame from "../../components/Frame.vue";
import InputField from "../../components/InputField.vue";
import Button from "../../components/Button.vue";
import User from "../../model/User";
import {AxiosError} from "axios";
import ServerError from "../../model/Error";

@Component({
    components: {Button, Frame, InputField}
})
export default class Passwords extends Vue {
    currentPassword = ""
    newPassword = ""
    repeatedPassword = ""

    showError = false
    showPasswordChangeSuccess = false

    errorMessage = ""

    updatePassword() {
        this.showPasswordChangeSuccess = false
        if (this.newPassword !== this.repeatedPassword) {
            this.errorMessage = "New password and repeated password fields don't correspond"
            this.showError = true
            return
        }
        let authenticatedUser: User | undefined = this.$store.state.auth.authentifiedUser;
        if (undefined === authenticatedUser) {
            this.errorMessage = "Could not detect user"
            this.showError = true
            return
        }
        this.$store.dispatch("auth/updatePassword", {
            id: authenticatedUser.id,
            currentPassword: this.currentPassword,
            newPassword: this.newPassword,
        })
            .then(() => {
                this.showError = false
                this.showPasswordChangeSuccess = true
                this.newPassword = ""
                this.repeatedPassword = ""
                this.currentPassword = ""
            })
            .catch((err: AxiosError<ServerError>) => {
                if (err.response === undefined) {
                    this.errorMessage = "Server error"
                } else {
                    this.errorMessage = "Error when changing password: " + err.response.data.error
                }
                this.showError = true
                this.showPasswordChangeSuccess = false
                console.error(err)
            })
    }
}
</script>

<style lang="scss">
.passwords {
    &__form {
        text-align: right;
    }
    &__title {
        margin-bottom: 1rem;
    }
    &__field {
        margin-top: 0.5rem;
        margin-bottom: 0.5rem;
        &-message {
            padding: 0.3rem 1rem;
            &--success {
                text-align: left;
                border: 1px solid $color-green;
                background-color: $color-green-light;
            }
            &--error {
                text-align: left;
                border: 1px solid $color-red;
                background-color: $color-red-light;
            }
        }
    }
}
</style>
