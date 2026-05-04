import { useEffect, useRef, useState, useMemo } from 'react';
import Pravah from './index';

/**
 * usePravah Hook
 * 
 * Initializes and manages a Pravah SDK instance.
 * 
 * @param {Object} options - Pravah initialization options { baseUrl, streams }
 * @returns {Pravah|null} The Pravah SDK instance
 */
export const usePravah = (options) => {
    const pravahRef = useRef(null);
    const optionsKey = useMemo(() => JSON.stringify(options), [options]);

    useEffect(() => {
        const pravah = new Pravah(options);
        pravahRef.current = pravah;

        return () => {
            if (pravah.isLeader) {
                pravah.closeEventSource();
                pravah.channel.postMessage({ type: 'LEADER_CLOSED' });
            }
            pravah.channel.close();
            pravahRef.current = null;
        };
    }, [optionsKey]);

    return pravahRef.current;
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
    useEffect(() => {
        if (!pravah || !event || !callback) return;

        pravah.on(event, callback);
        
        return () => {
            pravah.off(event, callback);
        };
    }, [pravah, event, callback]);
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
