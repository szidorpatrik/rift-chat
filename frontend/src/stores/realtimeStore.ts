import { defineStore } from 'pinia'
import { ref, shallowRef } from 'vue'

export const useRealtimeStore = defineStore('realtime', () => {
    // --- State ---
    const transport = shallowRef<WebTransport | null>(null)
    const isConnected = ref(false)

    // --- Actions ---
    async function connect() {
        if (transport.value) return

        try {
            const wtUrl = `https://localhost:8443/wt`

            transport.value = new WebTransport(wtUrl);

            await transport.value.ready
            isConnected.value = true
            console.log('WebTransport connected successfully.')

            // Start background listeners
            listenForClosed()
            listenForIncomingStreams()

        } catch (error) {
            console.error('WebTransport connection failed:', error)
            isConnected.value = false
        }
    }

    async function listenForClosed() {
        if (!transport.value) return
        try {
            await transport.value.closed
            console.log('WebTransport connection closed normally.')
        } catch (error) {
            console.error('WebTransport connection closed abruptly:', error)
        } finally {
            isConnected.value = false
            transport.value = null
        }
    }

    async function listenForIncomingStreams() {
        if (!transport.value) return

        const reader = transport.value.incomingUnidirectionalStreams.getReader()
        try {
            while (true) {
                const { done, value: stream } = await reader.read()
                if (done) break

                readStream(stream)
            }
        } catch (error) {
            console.error('Error reading incoming streams:', error)
        }
    }

    async function readStream(stream: ReadableStream) {
        const reader = stream.pipeThrough(new TextDecoderStream()).getReader()
        let buffer = ''

        try {
            while (true) {
                const { done, value } = await reader.read()
                if (done) break

                buffer += value

                const lines = buffer.split('\n')
                buffer = lines.pop() || ''

                for (const line of lines) {
                    if (line.trim()) {
                        const event = JSON.parse(line)
                        handleEvent(event)
                    }
                }
            }
        } catch (error) {
            console.error('Error reading stream data:', error)
        }
    }

    function handleEvent(event: any) {
        console.log('Received real-time event:', event)

        switch (event.type) {
            case 'CHANNEL_DELETED':
                // e.g., channelStore.removeChannel(event.channel_id)
                break;
            case 'NEW_MESSAGE':
                // e.g., messageStore.addMessage(event.message)
                break;
            default:
                console.warn('Unknown event type:', event.type)
        }
    }

    // --- Expose ---
    return {
        transport,
        isConnected,
        connect,
        listenForClosed,
        listenForIncomingStreams,
        readStream,
        handleEvent
    }
})