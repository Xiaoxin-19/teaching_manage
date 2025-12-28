import type { OrderTag } from '../types/appModels'

/**
 * 课时记录智能多标签分类算法 (v3.0)
 * 返回一个标签数组，支持复合场景识别
 */
export const categorizeOrderTags = (comment: string): OrderTag[] => {
  if (!comment) return [];
  const str = String(comment).trim();
  const tags: OrderTag[] = [];

  // --- 关键词定义 ---
  const keywordsCorrection = ['退费', '退款', '打错', '记错', '多记', '少记', '修正', '调整', '误操作', '补录', '撤销', '退回'];
  const keywordsSettlement = ['结余', '遗留', '旧系统', '初始', '剩余', '结算', '交接', '盘点', '底'];
  const keywordsTransfer = ['转给', '转入', '转出', '互转', '转课', '划拨'];
  const keywordsGift = ['送', '赠', '免', '奖', '集赞', '优惠', '福利', '奖励'];
  const keywordsReferral = ['介绍', '推荐', '拉新', '老带新', '拼团', '合报', '一起报'];
  const keywordsExpand = ['扩科', '新报', '加报', '多科'];
  // 支付方式关键词
  const keywordsPayment = ['微信', '支付宝', '现金', '银行', '扫码', '转账', '支付'];
  // 充值动词
  const keywordsRechargeAction = ['充值', '续费', '缴费', '补交', '买', '报名', '报课'];

  // --- 状态检测 ---
  const isCorrection = keywordsCorrection.some(k => str.includes(k));
  const isSettlement = keywordsSettlement.some(k => str.includes(k));
  const isTransfer = keywordsTransfer.some(k => str.includes(k));
  const hasGift = keywordsGift.some(k => str.includes(k));
  const hasReferral = keywordsReferral.some(k => str.includes(k));
  const hasExpand = keywordsExpand.some(k => str.includes(k));
  const hasPayment = keywordsPayment.some(k => str.includes(k));
  const hasRechargeAction = keywordsRechargeAction.some(k => str.includes(k));

  const isRecharge = hasPayment || hasRechargeAction;

  // --- 标签生成逻辑 ---

  // 1. 特殊/高优先级分类
  if (isCorrection) {
    tags.push({ label: '系统调整/纠错', color: 'error' });
    return tags;
  }
  if (isSettlement) {
    tags.push({ label: '历史结余', color: 'blue-grey' });
    return tags;
  }
  if (isTransfer) {
    tags.push({ label: '课时互转', color: 'indigo' });
  }

  // 2. 业务属性叠加
  if (hasReferral) {
    tags.push({ label: '转介绍/拼团', color: 'purple' });
  }
  if (hasExpand) {
    tags.push({ label: '扩科报名', color: 'deep-purple' });
  }
  if (hasPayment) {
    const paymentMethod = keywordsPayment.find(k => str.includes(k));
    if (paymentMethod) {
      tags.push({ label: `${paymentMethod}支付`, color: 'cyan-darken-1' });
    }
  }

  // 3. 核心业务定性
  if (isRecharge && hasGift) {
    tags.push({ label: '充值赠送', color: 'teal' });
  } else if (hasGift) {
    tags.push({ label: '活动赠送', color: 'orange' });
  } else if (isRecharge) {
    if (!hasExpand && !hasReferral) {
      tags.push({ label: '常规充值', color: 'success' });
    }
  }

  return tags;
};