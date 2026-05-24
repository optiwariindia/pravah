import React, { useEffect, useState, useMemo, createContext, useContext } from 'react';
import Pravah from './index';

const PravahContext = createContext(null);

/**
 * PravahProvider Component
 * 
 * Provides a managed Pravah instance to the component tree.
 * 
 * @param {Object} props
 * @param {Object} props.options - Pravah initialization options
 * @param {React.ReactNode} props.children
 */
export const PravahProvider = ({ options, children }) => {
    const pravah = usePravah(options);
    return (
        <PravahContext.Provider value={pravah}>
            {children}
        </PravahContext.Provider>
    );
};

/**
 * usePravahContext Hook
 * 
 * Accesses the Pravah instance from the context.
 */
export const usePravahContext = () => useContext(PravahContext);

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
    const optionsKey = useMemo(() => JSON.stringify(options), [options]);

    useEffect(() => {
        if (!options) {
            setPravah(null);
            return;
        }
        
        const instance = new Pravah(options);
        setPravah(instance);

        return () => {
            if (instance.isLeader) {
                instance.closeEventSource();
                instance.channel.postMessage({ type: 'LEADER_CLOSED' });
            }
            instance.channel.close();
            setPravah(null);
        };
    }, [optionsKey]);

    return pravah;
};

/**
 * usePravahListener Hook
 * 
 * Subscribes to a Pravah event and triggers a callback.
 * 
 * @param {Pravah|null} pravah - The Pravah SDK instance (optional if inside Provider)
 * @param {string} event - Event name to listen for
 * @param {Function} callback - Callback function
 */
export const usePravahListener = (pravah, event, callback) => {
    const contextPravah = usePravahContext();
    const activePravah = pravah || contextPravah;

    useEffect(() => {
        if (!activePravah || !event || !callback) return;

        activePravah.on(event, callback);
        
        return () => {
            activePravah.off(event, callback);
        };
    }, [activePravah, event, callback]);
};

/**
 * usePravahState Hook
 * 
 * A convenience hook that tracks the latest data for a specific event.
 * 
 * @param {Pravah|null} pravah - The Pravah SDK instance (optional if inside Provider)
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
