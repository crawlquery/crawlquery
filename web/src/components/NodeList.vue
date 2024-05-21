<script lang="ts" setup>
import { ref } from 'vue'
import moment from 'moment'
import axios from 'axios'
import { useAuthStore } from '../stores/auth'
import { Node } from '../types'
defineProps<{ nodes: Array<Node> }>()
const err = ref('')

const showKey = ref(false)
</script>
<template>

    <div class="overflow-x-auto">
        <table class="table table-zebra">
            <!-- head -->
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Key</th>
                    <th>Address</th>
                    <th>Shard ID</th>
                    <th>Created</th>
                </tr>
            </thead>
            <tbody>
                <!-- row 1 -->
                <tr v-for="node in nodes">
                    <td>{{ node.id }}</td>
                    <td>
                        <button @click="showKey = true" v-if="!showKey" class="btn btn-xs btn-info">Reveal key</button>
                        <span v-if="showKey">{{ node.key }}</span>
                    </td>
                    <td>{{ node.hostname }}:{{ node.port }}</td>
                    <td>{{ node.shard_id }}</td>
                    <td>{{ moment(node.created_at).fromNow() }}</td>
                </tr>

                <tr v-if="nodes.length === 0">
                    <td colspan="5" class="text-center">No nodes found</td>
                </tr>
            </tbody>
        </table>
    </div>
</template>