<template>
  <div class="navigation">
    <div class="navigation__logo-container">
      <Logo class="navigation__logo" />
    </div>
    <div class="navigation__menu">
      <div class="navigation__menu-main">
        <MenuItem v-for="link in links" :key="link.title" :title="link.title" @click="link.action"/>
      </div>
      <div class="navigation__menu-user">
        <MenuItem v-for="link in userLinks" :key="link.title" :title="link.title" @click="link.action" is-bottom="true"/>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Logo from "../components/navigation/Logo.vue";
import MenuItem from "../components/navigation/MenuItem.vue";

@Component({
  components: {MenuItem, Logo }
})
export default class Navigation extends Vue {
  get links(): {title: string, action: () => void}[] {
    let router = this.$router
    return [
      {
        title: "Classes",
        action: function () {
          router.push("/classes")
        }
      }
    ];
  }
  get userLinks(): {title: string, action: () => void}[] {
    let store = this.$store
    let router = this.$router
    return [
      {
        title: "Réglages",
        action: function () {
          router.push("/settings")
        }
      },
      {
        title: "Déconnexion",
        action: function () {
          store.dispatch("auth/signOut")
                  .catch((err: Error) => {
                    console.log("Sign out: Failure ("+err+")")
                  })
        }
      }
    ];
  }
}
</script>

<style lang="scss">
.navigation {
  display: flex;
  flex-direction: column;
}
.navigation {
  background-color: $color-primary-darker;
  position: relative;

  &__logo-container {
    background-color: $color-primary;
    display: flex;
    align-items: baseline;
    padding: 0.6rem;
  }

  &__logo {
    margin-right: 3.7rem;
    margin-left: 0.75rem;
  }

  &__profile {
    position: absolute;
    bottom: 20px;
    border-top: 1px solid $color-grey-7;
    color: $color-grey-8;
  }

  &__menu {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    flex-grow: 2;
  }
}
</style>
