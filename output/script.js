document.addEventListener('DOMContentLoaded', function() {
  const body = document.body;
  let classNames = [];
  classNames.push('device-iphone-xr');
  classNames.push('device-iphone-se');
  classNames.push('device-iphone-12-pro');
  classNames.push('device-ipad');
  classNames.push('device-iphone-pro-max');
  classNames.push('device-pixel-7');
  body.className = classNames.join(' ') + ' ' + body.className;
});
