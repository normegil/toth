<template>
<div class="authenticationscreen__root" v-if="!authenticated">
    <div class="authenticationscreen__dialog">
        <Logo class="authenticationscreen__logo" />
        <Login />
    </div>
</div>
</template>

<script lang="ts">
import Vue from "vue";
import Login from "./Login.vue";
import Component from "vue-class-component";
import Logo from "../../assets/images/logo.svg";

@Component({
    components: {Login, Logo}
})
export default class AuthenticationScreen extends Vue {
    get authenticated():boolean {
        let user = this.$store.state.auth.authentifiedUser;
        let b = undefined !== user && user.name !== "anonymous";
        return b
    }
}
</script>

<style lang="scss">
.authenticationscreen {
    &__root {
        position: absolute;
        top: 0;
        left: 0;
        min-width: 100vw;
        min-height: 100vh;

        z-index: 1;

        background-color: $color-primary-lighter;
    }

    &__dialog {
        @include center-block;
        margin: auto;
        padding: 2rem;
        background-color: $color-white;
        border: $color-grey-7 solid 2px;
    }

    &__logo {
        display: block;
        margin-left: auto;
        margin-right: auto;

        padding-bottom: 2rem;
        width: 15rem;
    }
}
</style>