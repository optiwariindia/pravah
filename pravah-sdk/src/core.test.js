import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import Pravah, { CONNECTION_STATES } from './core';

class FakeBroadcastChannel {
    static instances = [];

    constructor(name) {
        this.name = name;
        this.onmessage = null;
        this.closed = false;
        FakeBroadcastChannel.instances.push(this);
    }

    postMessage(data) {
        for (const channel of FakeBroadcastChannel.instances) {
            if (channel !== this && channel.name === this.name && !channel.closed) {
                channel.onmessage?.({ data });
            }
        }
    }

    close() {
        this.closed = true;
    }
}

class FakeEventSource {
    static instances = [];

    constructor(url, options) {
        this.url = url;
        this.options = options;
        this.listeners = {};
        this.closed = false;
        FakeEventSource.instances.push(this);
    }

    addEventListener(event, callback) {
        this.listeners[event] = callback;
    }

    removeEventListener(event) {
        delete this.listeners[event];
    }

    close() {
        this.closed = true;
    }

    dispatch(event, data = {}, lastEventId = '') {
        this.listeners[event]?.({ data: JSON.stringify(data), lastEventId });
    }
}

describe('Pravah core lifecycle', () => {
    beforeEach(() => {
        vi.useFakeTimers();
        FakeBroadcastChannel.instances = [];
        FakeEventSource.instances = [];
        vi.stubGlobal('BroadcastChannel', FakeBroadcastChannel);
        vi.stubGlobal('EventSource', FakeEventSource);
        vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
            ok: true,
            json: async () => ({ success: true, pravah_id: 'browser-1' })
        }));
        vi.stubGlobal('window', {
            addEventListener: vi.fn(),
            removeEventListener: vi.fn()
        });
    });

    afterEach(() => {
        vi.restoreAllMocks();
        vi.unstubAllGlobals();
        vi.useRealTimers();
    });

    it('elects one leader for simultaneous tabs', async () => {
        const options = { baseUrl: 'http://pravah', streams: ['project:1'], subscriptionToken: 'token-1' };
        const first = new Pravah(options);
        const second = new Pravah(options);

        await vi.advanceTimersByTimeAsync(900);

        expect(first.isLeader).not.toBe(second.isLeader);
        expect(FakeEventSource.instances.filter(source => !source.closed)).toHaveLength(1);
        first.close();
        second.close();
    });

    it('broadcasts leader signals to follower tabs', async () => {
        const options = { baseUrl: 'http://pravah', streams: ['project:1'], subscriptionToken: 'token-1' };
        const first = new Pravah(options);
        const second = new Pravah(options);
        const listener = vi.fn();
        const genericListener = vi.fn();
        second.on('bid-updated', listener);
        second.on('signal', genericListener);

        await vi.advanceTimersByTimeAsync(900);
        const source = FakeEventSource.instances.find(candidate => !candidate.closed);
        source.dispatch('signal', { event: 'bid-updated', data: { id: 'bid-1' } }, 'event-1');

        expect(listener).toHaveBeenCalledWith({ id: 'bid-1' });
        expect(genericListener).toHaveBeenCalledWith({
            event: 'bid-updated',
            data: { id: 'bid-1' }
        });
        expect(second.lastEventId).toBe('event-1');
        first.close();
        second.close();
    });

    it('elects a replacement when the active leader closes', async () => {
        const options = { baseUrl: 'http://pravah', streams: ['project:1'], subscriptionToken: 'token-1' };
        const first = new Pravah(options);
        const second = new Pravah(options);
        await vi.advanceTimersByTimeAsync(900);

        const leader = first.isLeader ? first : second;
        const follower = leader === first ? second : first;
        leader.close();
        await vi.advanceTimersByTimeAsync(350);

        expect(follower.isLeader).toBe(true);
        expect(FakeEventSource.instances.filter(source => !source.closed)).toHaveLength(1);
        follower.close();
    });

    it('does not share leadership across token scopes', async () => {
        const first = new Pravah({
            baseUrl: 'http://pravah',
            streams: ['project:1'],
            subscriptionToken: 'token-one'
        });
        const second = new Pravah({
            baseUrl: 'http://pravah',
            streams: ['project:1'],
            subscriptionToken: 'token-two'
        });
        await vi.advanceTimersByTimeAsync(900);

        expect(first.isLeader).toBe(true);
        expect(second.isLeader).toBe(true);
        first.close();
        second.close();
    });

    it('refreshes tokens on reconnect with bounded controlled retry', async () => {
        vi.stubGlobal('BroadcastChannel', undefined);
        const getSubscriptionToken = vi.fn()
            .mockResolvedValueOnce('token-1')
            .mockResolvedValue('token-2');
        const pravah = new Pravah({
            baseUrl: 'http://pravah',
            streams: ['project:1'],
            getSubscriptionToken
        });

        await vi.advanceTimersByTimeAsync(0);
        expect(FakeEventSource.instances[0].url).toContain('token=token-1');

        FakeEventSource.instances[0].onerror();
        await vi.advanceTimersByTimeAsync(1000);

        expect(getSubscriptionToken).toHaveBeenCalledTimes(2);
        expect(FakeEventSource.instances[1].url).toContain('token=token-2');
        pravah.close();
    });

    it('closes idempotently and cancels future reconnects', async () => {
        vi.stubGlobal('BroadcastChannel', undefined);
        const pravah = new Pravah({
            baseUrl: 'http://pravah',
            streams: ['project:1'],
            subscriptionToken: 'token-1'
        });
        await vi.advanceTimersByTimeAsync(0);
        FakeEventSource.instances[0].onerror();

        pravah.close();
        pravah.close();
        await vi.advanceTimersByTimeAsync(20000);

        expect(pravah.connectionState).toBe(CONNECTION_STATES.CLOSED);
        expect(FakeEventSource.instances).toHaveLength(1);
        expect(FakeEventSource.instances[0].closed).toBe(true);
    });

    it('reports missing EventSource without retrying', async () => {
        vi.stubGlobal('EventSource', undefined);
        const pravah = new Pravah({
            baseUrl: 'http://pravah',
            streams: ['project:1'],
            subscriptionToken: 'token-1'
        });
        const errorListener = vi.fn();
        pravah.on('error', errorListener);

        await vi.advanceTimersByTimeAsync(0);

        expect(pravah.connectionState).toBe(CONNECTION_STATES.DISCONNECTED);
        expect(errorListener).toHaveBeenCalledOnce();
        expect(FakeEventSource.instances).toHaveLength(0);
        pravah.close();
    });
});
