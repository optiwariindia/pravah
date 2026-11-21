import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook } from '@testing-library/react';
import { usePravah, usePravahListener } from './index';

// Mock the core class inside the same package
const mockOn = vi.fn();
const mockOff = vi.fn();
const mockCloseEventSource = vi.fn();
const mockChannelClose = vi.fn();
const mockPostMessage = vi.fn();

vi.mock('./core', () => {
    return {
        default: vi.fn().mockImplementation(() => ({
            on: mockOn,
            off: mockOff,
            close: mockCloseEventSource,
            isLeader: true
        }))
    };
});

describe('pravah-react hooks', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('usePravah should initialize the SDK', () => {
        const { result, unmount } = renderHook(() => usePravah({ baseUrl: 'http://test', streams: ['test'] }));
        expect(result.current).toBeDefined();
        expect(mockOn).toHaveBeenCalledWith('statechange', expect.any(Function));
        unmount();
        expect(mockOff).toHaveBeenCalledWith('statechange', expect.any(Function));
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
