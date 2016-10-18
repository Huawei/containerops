








var dom = {};




export function ContextMenu($span,items,opt) {
	
	
	menuHide();


  	// create root element
  	var root = document.createElement('div');
  	root.className = 'jsoneditor-contextmenu-root-treeedit';
  	dom.root = root;

  	// create a container element
  	var menu = document.createElement('div');
  	menu.className = 'jsoneditor-contextmenu';
  	dom.menu = menu;
  	root.appendChild(menu);

  	// create a list to hold the menu items
  	var list = document.createElement('ul');
  	list.className = 'jsoneditor-menu';
  	menu.appendChild(list);
  	dom.list = list;
  	dom.items = []; // list with all buttons


  	// create a (non-visible) button to set the focus to the menu
  	var focusButton = document.createElement('button');
  	dom.focusButton = focusButton;
  	var li = document.createElement('li');
  	li.style.overflow = 'hidden';
  	li.style.height = '0';
  	li.appendChild(focusButton);
  	list.appendChild(li);

  	
  	createMenuItems(list,dom.items,items,opt);

  	$span.before(root);

};


function createMenuItems (list, domItems,items,opt) {

    items.forEach(function (item) {
      if (item.type == 'separator') {
        // create a separator
        var separator = document.createElement('div');
        separator.className = 'jsoneditor-separator';
        var li = document.createElement('li');
        li.appendChild(separator);
        list.appendChild(li);
      }
      else {
        var domItem = {};

        // create a menu item
        var li = document.createElement('li');
        list.appendChild(li);

        var button = document.createElement('button');
        button.className = item.className;
        domItem.button = button;
        if (item.title) {
          button.title = item.title;
        }
        if (item.click) {
          button.onclick = function (event) {
            event.preventDefault();
            event.stopPropagation();
            item.click($(this),opt);
            menuHide();
            
          };
        }
        li.appendChild(button);

        if (item.submenu) {
        	// add the icon to the button
          	var divIcon = document.createElement('div');
          	divIcon.className = 'jsoneditor-icon';
          	button.appendChild(divIcon);
          	button.appendChild(document.createTextNode(item.text));
          	var buttonSubmenu;

          	if (item.click) {
	            // submenu and a button with a click handler
	            button.className += ' jsoneditor-default';

	            var buttonExpand = document.createElement('button');
	            domItem.buttonExpand = buttonExpand;
	            buttonExpand.className = 'jsoneditor-expand';
	            buttonExpand.innerHTML = '<div class="jsoneditor-expand"></div>';
	            li.appendChild(buttonExpand);
	            if (item.submenuTitle) {
	              buttonExpand.title = item.submenuTitle;
	            }

	            buttonSubmenu = buttonExpand;
	          }
	          else {
	            // submenu and a button without a click handler
	            var divExpand = document.createElement('div');
	            divExpand.className = 'jsoneditor-expand';
	            button.appendChild(divExpand);

	            buttonSubmenu = button;
	          }


	          // attach a handler to expand/collapse the submenu
	          buttonSubmenu.onclick = function (event) {
	            event.preventDefault();
	            event.stopPropagation();
	            let _ul = $(this).parents("li").find("ul");

	            if(_ul.is(":hidden")){
	            	_ul.height(_ul.find("li").length * 24).css("display","block");
	            }else{
					_ul.height(0).css("display","");
	            }
	          };

	          // create the submenu
	          var domSubItems = [];
	          domItem.subItems = domSubItems;
	          var ul = document.createElement('ul');
	          domItem.ul = ul;
	          ul.className = 'jsoneditor-menu';
	          ul.style.height = '0';
	          li.appendChild(ul);
	          createMenuItems(ul, domSubItems, item.submenu,opt);


        }else {
          // no submenu, just a button with clickhandler
          button.innerHTML = '<div class="jsoneditor-icon"></div>' + item.text;
        }


        domItems.push(domItem);

      }


    });


	$("body").click(function(){
		$(".jsoneditor-contextmenu-root-treeedit").remove();
	});
}

function menuHide(){
	$(".jsoneditor-contextmenu-root-treeedit").remove();
}


