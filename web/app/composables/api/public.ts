export type PluginListQuery = {
    type?: number
    frameId?: number
    tagId?: number
}

const buildPluginListQuery = (page: number, pageSize: number, query: PluginListQuery = {}) => ({
    page,
    pageSize,
    ...(query.type !== undefined && query.type !== -1 ? { type: query.type } : {}),
    ...(query.frameId !== undefined && query.frameId !== -1 ? { frameId: query.frameId } : {}),
    ...(query.tagId !== undefined && query.tagId !== -1 ? { tagId: query.tagId } : {}),
})

export const getHotPlugins = (pageSize: number = 10, page: number = 1, query: PluginListQuery = {}) => {
    return useApiFetch('/plugin/list/hot', {
        method: 'GET',
        query: buildPluginListQuery(page, pageSize, query)
    })
}

export const getLatestPlugins = (pageSize: number = 10, page: number = 1, query: PluginListQuery = {}) => {
    return useApiFetch('/plugin/list/latest', {
        method: 'GET',
        query: buildPluginListQuery(page, pageSize, query)
    })
}

export const searchPlugins = (keywords: string, pageSize: number = 20, page: number = 1, query: PluginListQuery = {}) => {
    return useApiFetch('/plugin/search', {
        method: 'GET',
        query: { keywords, ...buildPluginListQuery(page, pageSize, query) }
    })
}

export const getPluginsByUser = (userId: string | number, pageSize: number = 20, page: number = 1, query: PluginListQuery = {}) => {
    return useApiFetch('/plugin/list/by-user', {
        method: 'GET',
        query: { userId, ...buildPluginListQuery(page, pageSize, query) }
    })
}

export const getStarPlugins = (pageSize: number = 20, page: number = 1) => {
    return useApiFetch('/star/list', {
        method: 'GET',
        query: { page, pageSize }
    })
}

export const toggleStar = (pluginId: number) => {
    return useApiFetch<{ is_starred: boolean }>('/star/toggle', {
        method: 'POST',
        body: { pluginId },
    })
}

export const getMyPlugins = (pageSize: number = 10, page: number = 1, query: Record<string, any> = {}) => {
    return useApiFetch('/plugin/my', {
        method: 'GET',
        query: { page, pageSize, ...query }
    })
}

export const getFrameList = (pageSize: number = 100, page: number = 1, keywords: string = '') => {
    return useApiFetch('/frame/list', { method: 'GET', query: { page, pageSize, ...(keywords ? { keywords } : {}) } })
}

export const getTagList = (pageSize: number = 100, page: number = 1, keywords: string = '') => {
    return useApiFetch('/tag/list', { method: 'GET', query: { page, pageSize, ...(keywords ? { keywords } : {}) } })
}
