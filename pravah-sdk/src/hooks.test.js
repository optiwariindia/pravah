import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook } from '@testing-library/react';
import { usePravah, usePravahListener } from './index';

// Mock the core class inside the same package
const mockOn = vi.fn();
const mockOff = vi.fn();
const mockCloseEventSource = vi.fn();
const mockChannelClose = vi.fn();
const mockPostMessage = vi.fn();

vi.mock('./index', async (importOriginal) => {
    const actual = await importOriginal();
    return {
        ...actual,
        default: vi.fn().mockImplementation(() => {
            return {
                on: mockOn,
                off: mockOff,
                closeEventSource: mockCloseEventSource,
                channel: {
                    close: mockChannelClose,
                    postMessage: mockPostMessage
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

    it('usePravah should initialize the SDK', () => {
        const { result } = renderHook(() => usePravah({ baseUrl: 'http://test', streams: ['test'] }));
        expect(result.current).toBeDefined();
        expect(mockOn).not.toHaveBeenCalled(); // Handled by listeners
    });

    it('usePravahListener should register and unregister listeners', () => {
        const mockPravah = { on: mockOn, off: mockOff };
        const callback = vi.fn();
        
        const { unmount } = renderHook(() => 
            usePravahListener(mockPravah, 'test-event', callback)
        );

        expect(mockOn).toHaveBeenCalledWith('test-event', callback);
        
        unmount();
        
        expect(mockOff).toHaveBeenCalledWith('test-event', callback);
    });
});
