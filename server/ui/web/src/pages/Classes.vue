<template>
    <div>
        <Frame>
            <h1>Classes</h1>
        </Frame>
        <Frame>
            <table>
                <tr>
                    <th>
                        Nom
                    </th>
                    <th>
                        Titulaire
                    </th>
                    <th>
                        Nombre d'élèves
                    </th>
                </tr>
            </table>
        </Frame>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Passwords from "../layout/settings/Passwords.vue";
import ProfileSettings from "../layout/settings/ProfileSettings.vue";
import Frame from "../components/Frame.vue";

@Component({
    components: {Frame}
})
export default class Classes extends Vue {
}
</script>

<style lang="scss">
</style>
