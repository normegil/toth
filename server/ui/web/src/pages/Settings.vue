<template>
    <div>
        <Frame>
            <h1>Réglages</h1>
        </Frame>
        <div class="settings__container-profile">
            <ProfileSettings class="settings__subcontainer-profile" />
            <Passwords class="settings__subcontainer-profile" />
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Passwords from "../layout/settings/Passwords.vue";
import ProfileSettings from "../layout/settings/ProfileSettings.vue";
import Frame from "../components/Frame.vue";

@Component({
    components: {ProfileSettings, Passwords, Frame}
})
export default class Settings extends Vue {
}
</script>

<style lang="scss">
.frame {
    margin-bottom: 2rem;
}
.settings {
    &__container {
        &-profile {
            display: flex;
        }
    }
    &__subcontainer {
        &-profile {
            flex-grow: 1;

            &:first-child {
                margin-right: 1rem;
            }
            &:last-child {
                margin-left: 1rem;
            }
        }
    }
}
</style>
