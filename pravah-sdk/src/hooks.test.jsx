import React from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { usePravah, usePravahListener, PravahProvider } from './hooks';
import Pravah from './index';

// Mock only the default export of index.js (the Pravah class)
vi.mock('./index', () => {
    const mockOn = vi.fn();
    const mockOff = vi.fn();
    return {
        default: vi.fn().mockImplementation(() => {
            return {
                on: mockOn,
                off: mockOff,
                closeEventSource: vi.fn(),
                channel: {
                    close: vi.fn(),
                    postMessage: vi.fn()
                },
                isLeader: true
            };
        })
    };
});

describe('pravah-react hooks', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('usePravah should initialize the SDK', async () => {
        const { result } = renderHook(() => usePravah({ baseUrl: 'http://test', streams: ['test'] }));
        await waitFor(() => {
            expect(result.current).toBeDefined();
        });
    });

    it('usePravahListener should register and unregister listeners', () => {
        const mockPravah = new Pravah();
        const callback = vi.fn();
        
        const { unmount } = renderHook(() => 
            usePravahListener(mockPravah, 'test-event', callback)
        );

        expect(mockPravah.on).toHaveBeenCalledWith('test-event', callback);
        
        unmount();
        
        expect(mockPravah.off).toHaveBeenCalledWith('test-event', callback);
    });

    it('usePravahListener should use pravah from context if not provided', async () => {
        const callback = vi.fn();
        
        const { unmount } = renderHook(() => usePravahListener(null, 'context-event', callback), {
            wrapper: ({ children }) => (
                <PravahProvider options={{ baseUrl: 'http://test' }}>
                    {children}
                </PravahProvider>
            )
        });

        // We need to get the instance created by the Provider to check calls
        // Since we mocked the constructor, we can check the calls on the mock instances
        const mockInstance = vi.mocked(Pravah).mock.results[0].value;

        await waitFor(() => {
            expect(mockInstance.on).toHaveBeenCalledWith('context-event', callback);
        });

        unmount();
        expect(mockInstance.off).toHaveBeenCalledWith('context-event', callback);
    });
});
