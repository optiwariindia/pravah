import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import Pravah from './index';

describe('Pravah SDK logLevel', () => {
    let logSpy, warnSpy, errorSpy;

    beforeEach(() => {
        logSpy = vi.spyOn(console, 'log').mockImplementation(() => {});
        warnSpy = vi.spyOn(console, 'warn').mockImplementation(() => {});
        errorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    });

    afterEach(() => {
        logSpy.mockRestore();
        warnSpy.mockRestore();
        errorSpy.mockRestore();
    });

    it('should not log anything by default (logLevel = none)', () => {
        const pravah = new Pravah({ logLevel: 'none' });
        pravah._log('info', 'test info');
        pravah._log('warn', 'test warn');
        pravah._log('error', 'test error');

        expect(logSpy).not.toHaveBeenCalled();
        expect(warnSpy).not.toHaveBeenCalled();
        expect(errorSpy).not.toHaveBeenCalled();
    });

    it('should log only errors when logLevel is error', () => {
        const pravah = new Pravah({ logLevel: 'error' });
        pravah._log('info', 'test info');
        pravah._log('warn', 'test warn');
        pravah._log('error', 'test error');

        expect(logSpy).not.toHaveBeenCalled();
        expect(warnSpy).not.toHaveBeenCalled();
        expect(errorSpy).toHaveBeenCalledWith('test error');
    });

    it('should log errors and warnings when logLevel is warn', () => {
        const pravah = new Pravah({ logLevel: 'warn' });
        pravah._log('info', 'test info');
        pravah._log('warn', 'test warn');
        pravah._log('error', 'test error');

        expect(logSpy).not.toHaveBeenCalled();
        expect(warnSpy).toHaveBeenCalledWith('test warn');
        expect(errorSpy).toHaveBeenCalledWith('test error');
    });

    it('should log everything when logLevel is info', () => {
        const pravah = new Pravah({ logLevel: 'info' });
        pravah._log('info', 'test info');
        pravah._log('warn', 'test warn');
        pravah._log('error', 'test error');

        expect(logSpy).toHaveBeenCalledWith('test info');
        expect(warnSpy).toHaveBeenCalledWith('test warn');
        expect(errorSpy).toHaveBeenCalledWith('test error');
    });

    it('should default to none if no logLevel is provided', () => {
        const pravah = new Pravah();
        pravah._log('error', 'test error');
        expect(errorSpy).not.toHaveBeenCalled();
    });
});
