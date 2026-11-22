import { useEffect, useRef, useState } from 'react';
import Pravah from './core';

/**
 * usePravah Hook
 * 
 * Initializes and manages a Pravah SDK instance.
 * 
 * @param {Object} options - Pravah initialization options { baseUrl, streams }
 * @returns {Pravah|null} The Pravah SDK instance
 */
export const usePravah = (options) => {
    const [pravah, setPravah] = useState(null);
    const [, setConnectionRevision] = useState(0);
    const streamsKey = JSON.stringify(options?.streams || ['global']);

    useEffect(() => {
        const instance = new Pravah(options);
        const handleStateChange = () => setConnectionRevision(revision => revision + 1);
        instance.on('statechange', handleStateChange);
        setPravah(instance);

        return () => {
            instance.off('statechange', handleStateChange);
            instance.close();
            setPravah(null);
        };
    }, [
        options?.baseUrl,
        streamsKey,
        options?.subscriptionToken,
        options?.getSubscriptionToken
    ]);

    return pravah;
};

/**
 * usePravahListener Hook
 * 
 * Subscribes to a Pravah event and triggers a callback.
 * 
 * @param {Pravah} pravah - The Pravah SDK instance
 * @param {string} event - Event name to listen for
 * @param {Function} callback - Callback function
 */
export const usePravahListener = (pravah, event, callback) => {
    const callbackRef = useRef(callback);
    callbackRef.current = callback;

    useEffect(() => {
        if (!pravah || !event || !callback) return;

        const listener = (payload) => callbackRef.current?.(payload);
        pravah.on(event, listener);
        
        return () => {
            pravah.off(event, listener);
        };
    }, [pravah, event]);
};

/**
 * usePravahState Hook
 * 
 * A convenience hook that tracks the latest data for a specific event.
 * 
 * @param {Pravah} pravah - The Pravah SDK instance
 * @param {string} event - Event name to listen for
 * @returns {any} The latest event data
 */
export const usePravahState = (pravah, event) => {
    const [data, setData] = useState(null);

    usePravahListener(pravah, event, (payload) => {
        setData(payload);
    });

    return data;
};
