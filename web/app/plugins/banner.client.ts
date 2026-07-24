export default defineNuxtPlugin(() => {
    console.info(
        '%c PluginMarket %c 插件市场 %c\n PoweredBy %c 周星星 %c\n GitHub %c https://github.com/xiaoyanu/PluginMarket ',
        'margin: 1em 0;padding: 5px 0px; border-radius: 5px 0 0 5px; color: #fff; background: #FF6699; font-weight: bold;',
        'padding: 5px 0px; border-radius: 0 5px 5px 0; color: #fff; background: #FF9999; font-weight: bold;',
        'color: #fff;  padding: 5px 0; background: #3498db;border-radius: 5px 0 0 5px;',
        'margin-bottom: 1em; padding: 5px 0; background: #efefef;border-radius: 0 5px 5px 0;',
        'color: #fff;  padding: 5px 0; background: #000;border-radius: 5px 0 0 5px;',
        'margin-bottom: 1em; padding: 5px 0; background: #efefef;border-radius: 0 5px 5px 0;'
    )
})