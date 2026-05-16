import type { HostProvider, ProviderRegistry } from './types.ts';

class ProviderRegistryImpl implements ProviderRegistry {
  private providers: HostProvider[] = [];

  register(provider: HostProvider): void {
    // 检查重复 ID
    if (this.providers.some((p) => p.id === provider.id)) {
      throw new Error(`Provider with id "${provider.id}" already registered`);
    }
    this.providers.push(provider);
  }

  findProvider(url: string): HostProvider | null {
    for (const provider of this.providers) {
      const match = provider.match(url);
      if (match.matches) {
        return provider;
      }
    }
    return null;
  }

  getProviders(): HostProvider[] {
    return [...this.providers];
  }
}

// 单例注册表
export const registry = new ProviderRegistryImpl();

/** 向全局注册表注册 provider。 */
export function registerProvider(provider: HostProvider): void {
  registry.register(provider);
}

/** 查找匹配给定 URL 的 provider。 */
export function findProvider(url: string): HostProvider | null {
  return registry.findProvider(url);
}

/** 获取所有已注册的 provider。 */
export function getProviders(): HostProvider[] {
  return registry.getProviders();
}
