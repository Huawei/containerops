import {Directive, ElementRef} from '@angular/core';
declare var $: any;

@Directive ({
  selector: '[check-all]'
})

export class CheckAllDirective {
  $el: any;

  constructor(el: ElementRef) {
    this.$el = $(el.nativeElement);
  }

  ngOnInit(): void {
    let $el = this.$el;
    $el.on('click', function(): void {
      $el.closest('table').find('input[type=checkbox]')
        .not(this).prop('checked', $(this).prop('checked'));
    });
  }
}
