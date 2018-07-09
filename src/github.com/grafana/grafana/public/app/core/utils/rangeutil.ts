import _ from 'lodash';
import moment from 'moment';
import * as dateMath from './datemath';

var spans = {
  s: { display: '초' },
  m: { display: '분' },
  h: { display: '시' },
  d: { display: '일' },
  w: { display: '주' },
  M: { display: '월' },
  y: { display: '년' },
};

var rangeOptions = [
  { from: 'now/d', to: 'now/d', display: '오늘', section: 2 },
  { from: 'now/d', to: 'now', display: '현재까지', section: 2 },
  { from: 'now/w', to: 'now/w', display: '이번 주', section: 2 },
  { from: 'now/w', to: 'now', display: '이번주 지금까지', section: 2 },
  { from: 'now/M', to: 'now/M', display: '이번 달', section: 2 },
  { from: 'now/M', to: 'now', display: '이번 달 지금까지', section: 2 },
  { from: 'now/y', to: 'now/y', display: '올해', section: 2 },
  { from: 'now/y', to: 'now', display: '올해 지금까지', section: 2 },

  { from: 'now-1d/d', to: 'now-1d/d', display: '어제', section: 1 },
  {
    from: 'now-2d/d',
    to: 'now-2d/d',
    display: '어제 전날',
    section: 1,
  },
  {
    from: 'now-7d/d',
    to: 'now-7d/d',
    display: '지난 주 오늘',
    section: 1,
  },
  { from: 'now-1w/w', to: 'now-1w/w', display: '이전 주', section: 1 },
  { from: 'now-1M/M', to: 'now-1M/M', display: '이전 월', section: 1 },
  { from: 'now-1y/y', to: 'now-1y/y', display: '이전 해', section: 1 },

  { from: 'now-5m', to: 'now', display: '지난 5분', section: 3 },
  { from: 'now-15m', to: 'now', display: '지난 15분', section: 3 },
  { from: 'now-30m', to: 'now', display: '지난 30분', section: 3 },
  { from: 'now-1h', to: 'now', display: '지난 1시간', section: 3 },
  { from: 'now-3h', to: 'now', display: '지난 3시간', section: 3 },
  { from: 'now-6h', to: 'now', display: '지난 6시간', section: 3 },
  { from: 'now-12h', to: 'now', display: '지난 12시간', section: 3 },
  { from: 'now-24h', to: 'now', display: '지난 24시간', section: 3 },

  { from: 'now-2d', to: 'now', display: '이틀 전', section: 0 },
  { from: 'now-7d', to: 'now', display: '7일 전', section: 0 },
  { from: 'now-30d', to: 'now', display: '30일 전', section: 0 },
  { from: 'now-90d', to: 'now', display: '90일 전', section: 0 },
  { from: 'now-6M', to: 'now', display: '6개월 전', section: 0 },
  { from: 'now-1y', to: 'now', display: '1년 전', section: 0 },
  { from: 'now-2y', to: 'now', display: '2년 전', section: 0 },
  { from: 'now-5y', to: 'now', display: '5년 전', section: 0 },
];

var absoluteFormat = 'MMM D, YYYY HH:mm:ss';

var rangeIndex = {};
_.each(rangeOptions, function(frame) {
  rangeIndex[frame.from + ' to ' + frame.to] = frame;
});

export function getRelativeTimesList(timepickerSettings, currentDisplay) {
  var groups = _.groupBy(rangeOptions, (option: any) => {
    option.active = option.display === currentDisplay;
    return option.section;
  });

  // _.each(timepickerSettings.time_options, (duration: string) => {
  //   let info = describeTextRange(duration);
  //   if (info.section) {
  //     groups[info.section].push(info);
  //   }
  // });

  return groups;
}

function formatDate(date) {
  return date.format(absoluteFormat);
}

// handles expressions like
// 5m
// 5m to now/d
// now/d to now
// now/d
// if no to <expr> then to now is assumed
export function describeTextRange(expr: any) {
  let isLast = expr.indexOf('+') !== 0;
  if (expr.indexOf('now') === -1) {
    expr = (isLast ? 'now-' : 'now') + expr;
  }

  let opt = rangeIndex[expr + ' to now'];
  if (opt) {
    return opt;
  }

  if (isLast) {
    opt = { from: expr, to: 'now' };
  } else {
    opt = { from: 'now', to: expr };
  }

  let parts = /^now([-+])(\d+)(\w)/.exec(expr);
  if (parts) {
    let unit = parts[3];
    let amount = parseInt(parts[2]);
    let span = spans[unit];
    if (span) {
      opt.display = isLast ? '지난 ' : '다음 ';
      opt.display += amount + ' ' + span.display;
      opt.section = span.section;
      if (amount > 1) {
        opt.display += 's';
      }
    }
  } else {
    opt.display = opt.from + ' to ' + opt.to;
    opt.invalid = true;
  }

  return opt;
}

export function describeTimeRange(range) {
  var option = rangeIndex[range.from.toString() + ' to ' + range.to.toString()];
  if (option) {
    return option.display;
  }

  if (moment.isMoment(range.from) && moment.isMoment(range.to)) {
    return formatDate(range.from) + ' to ' + formatDate(range.to);
  }

  if (moment.isMoment(range.from)) {
    var toMoment = dateMath.parse(range.to, true);
    return formatDate(range.from) + ' to ' + toMoment.fromNow();
  }

  if (moment.isMoment(range.to)) {
    var from = dateMath.parse(range.from, false);
    return from.fromNow() + ' to ' + formatDate(range.to);
  }

  if (range.to.toString() === 'now') {
    var res = describeTextRange(range.from);
    return res.display;
  }

  return range.from.toString() + ' to ' + range.to.toString();
}
