package svc

import (
	"context"

	"gitee.com/chunanyong/zorm"
	gtxContext "github.com/cectc/hptx/pkg/base/context"
	"github.com/cectc/hptx/pkg/tm"
)

// 建议以下代码放到单独的文件里
//................//

// ZormGlobalTransaction 包装seata/hptx的*tm.DefaultGlobalTransaction,实现zorm.IGlobalTransaction接口
type ZormGlobalTransaction struct {
	*tm.DefaultGlobalTransaction
}

// MyFuncGlobalTransaction zorm适配seata/hptx 全局分布式事务的函数
// 重要!!!!需要配置zorm.DataSourceConfig.FuncGlobalTransaction=MyFuncGlobalTransaction 重要!!!
func MyFuncGlobalTransaction(ctx context.Context) (zorm.IGlobalTransaction, context.Context, error) {
	//获取seata/hptx的rootContext
	rootContext := gtxContext.NewRootContext(ctx)
	//创建seata/hptx事务
	globalTx := tm.GetCurrentOrCreate(rootContext)
	//使用zorm.IGlobalTransaction接口对象包装分布式事务,隔离seata/hptx依赖
	globalTransaction := &ZormGlobalTransaction{globalTx}

	return globalTransaction, rootContext, nil
}

//实现zorm.IGlobalTransaction 托管全局分布式事务接口,seata和hptx目前实现代码一致,只是引用的实现包不同
// Begin 开启全局分布式事务
func (gtx *ZormGlobalTransaction) BeginGTX(ctx context.Context, globalRootContext context.Context) error {
	rootContext := globalRootContext.(*gtxContext.RootContext)
	return gtx.BeginWithTimeout(int32(6000), rootContext)
}

// Commit 提交全局分布式事务
func (gtx *ZormGlobalTransaction) CommitGTX(ctx context.Context, globalRootContext context.Context) error {
	rootContext := globalRootContext.(*gtxContext.RootContext)
	return gtx.Commit(rootContext)
}

// Rollback 回滚全局分布式事务
func (gtx *ZormGlobalTransaction) RollbackGTX(ctx context.Context, globalRootContext context.Context) error {
	rootContext := globalRootContext.(*gtxContext.RootContext)
	//如果是Participant角色,修改为Launcher角色,允许分支事务提交全局事务.
	if gtx.Role != tm.Launcher {
		gtx.Role = tm.Launcher
	}
	return gtx.Rollback(rootContext)
}

// GetXID 获取全局分布式事务的XID
func (gtx *ZormGlobalTransaction) GetGTXID(ctx context.Context, globalRootContext context.Context) string {
	rootContext := globalRootContext.(*gtxContext.RootContext)
	return rootContext.GetXID()
}

//................//
