<template>
    <Frame class="profile-settings__container">
        <h2 class="profile-settings__title">Profil</h2>
        <div class="profile-settings__form">
            <InputField class="profile-settings__field" label="Nom" v-model="name" />
            <InputField class="profile-settings__field" label="Mail" v-model="mail" />
            <Button title="Enregistrer" @click="updateProfile" />
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

@Component({
    components: {Frame, InputField, Button}
})
export default class ProfileSettings extends Vue {
    name:string = ""
    mail:string = ""

    get authentifiedUser(): User {
        return this.$store.state.auth.authentifiedUser
    }

    mounted() {
        this.name = this.authentifiedUser.name
        this.mail = this.authentifiedUser.mail
    }

    updateProfile() {
        this.$store.dispatch("auth/updateProfile", {
            id: this.authentifiedUser.id,
            name: this.name,
            mail: this.mail,
        })
            .catch((err: Error) => {
                console.error(err)
            })
    }
}
</script>

<style lang="scss">
.profile-settings {
    &__form {
        text-align: right;
    }
    &__title {
        margin-bottom: 1rem;
    }
    &__field {
        margin-top: 0.5rem;
        margin-bottom: 0.5rem;
    }
}
</style>
