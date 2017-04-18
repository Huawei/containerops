import { Angular4Page } from './app.po';

describe('angular4 App', () => {
  let page: Angular4Page;

  beforeEach(() => {
    page = new Angular4Page();
  });

  it('should display message saying app works', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('app works!');
  });
});
