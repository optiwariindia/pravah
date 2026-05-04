import { useEffect, useRef, useState } from 'react';

/**
 * usePravah Hook
 * 
 * Provides a React-idiomatic way to interact with the Pravah SDK.
 * Handles instance lifecycle, connection status, and event subscription.
 * 
 * @param {Object} options - Pravah initialization options { baseUrl, streams }
 * @returns {Object} { pravah, isConnected, lastSignal }
 */
export const usePravah = (options) => {
    const [isConnected, setIsConnected] = useState(false);
    const [lastSignal, setLastSignal] = useState(null);
    const pravahRef = useRef(null);

    // Ensure options are stable or handled correctly to prevent infinite re-inits
    const optionsKey = JSON.stringify(options);

    useEffect(() => {
        // Initialize the SDK instance
        // Assuming Pravah is available globally or imported
        if (typeof Pravah === 'undefined') {
            console.error('[usePravah] Pravah SDK not found. Please ensure pravah-sdk.js is loaded.');
            return;
        }

        const pravah = new Pravah(options);
        pravahRef.current = pravah;

        // Listen for internal connection status
        pravah.on('connected', (data) => {
            setIsConnected(true);
            setLastSignal({ event: 'connected', data });
        });

        pravah.on('disconnected', () => {
            setIsConnected(false);
        });

        // Cleanup on unmount
        return () => {
            if (pravah.isLeader) {
                pravah.closeEventSource();
            }
            // Notify other tabs if this was the leader
            if (pravah.isLeader) {
                pravah.channel.postMessage({ type: 'LEADER_CLOSED' });
            }
            pravah.channel.close();
        };
    }, [optionsKey]);

    /**
     * registerListener
     * Helper to subscribe to specific signals within a component
     */
    const registerListener = (event, callback) => {
        useEffect(() => {
            if (!pravahRef.current) return;

            const internalCallback = (data) => {
                setLastSignal({ event, data });
                callback(data);
            };

            pravahRef.current.on(event, internalCallback);

            // Note: Currently the base SDK doesn't support 'off' (unsubscribing).
            // This is a great candidate for a community feature request!
        }, [event, callback]);
    };

    return {
        pravah: pravahRef.current,
        isConnected,
        lastSignal,
        registerListener
    };
};
