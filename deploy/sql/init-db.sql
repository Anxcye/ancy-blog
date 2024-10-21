
-- create database ancy_blog if not exists;

use ancy_blog;

create table ancy_article
(
    id          bigint auto_increment
        primary key,
    title       varchar(256)        null comment '标题',
    content     longtext            null comment '文章内容',
    summary     varchar(1024)       null comment '文章摘要',
    category_id bigint  default 0   null comment '所属分类id',
    thumbnail   varchar(256)        null comment '缩略图',
    is_top      char    default '0' null comment '是否置顶（0否，1是）',
    status      char    default '1' null comment '状态（0已发布，1草稿）',
    type        tinyint default 0   not null comment '文章类型 0普通文章 1首页文章 2友链文章',
    order_num   tinyint default 0   not null comment '排序',
    view_count  bigint  default 0   null comment '访问量',
    is_comment  char    default '1' null comment '是否允许评论 1是，0否',
    create_by   bigint              null,
    create_time datetime            null,
    update_by   bigint              null,
    update_time datetime            null,
    deleted     int     default 0   null comment '删除标志（0代表未删除，1代表已删除）'
)
    comment '文章表';



create table ancy_article_tag
(
    id         bigint auto_increment comment 'ID'
        primary key,
    article_id bigint           not null comment '文章id',
    tag_id     bigint default 0 not null comment '标签id'
)
    comment '文章标签关联表';

create table ancy_category
(
    id          bigint auto_increment
        primary key,
    name        varchar(128)       null comment '分类名',
    parent_id   bigint default -1  null comment '父分类id，如果没有父分类为-1',
    description varchar(512)       null comment '描述',
    status      char   default '0' null comment '状态0:正常,1禁用',
    create_by   bigint             null,
    create_time datetime           null,
    update_by   bigint             null,
    update_time datetime           null,
    deleted     int    default 0   null comment '删除标志（0代表未删除，1代表已删除）'
)
    comment '分类表';

create table ancy_comment
(
    id                  bigint auto_increment
        primary key,
    type                char   default '0' null comment '评论类型（0代表文章评论，1代表NOTE评论）',
    article_id          bigint             null comment '文章id',
    status              char   default '0' null comment '公开状态 0代表公开，1代表隐藏',
    parent_id           bigint default -1  not null comment '根评论id',
    user_id             bigint default 0   null,
    avatar              varchar(300)       null,
    nickname            varchar(30)        null,
    email               varchar(255)       null,
    content             varchar(1024)      null comment '评论内容',
    ua                  varchar(300)       null,
    ip                  varchar(50)        null,
    like_count          int    default 0   null,
    is_top              char   default '0' null,
    to_comment_nickname varchar(30)        null,
    to_comment_id       bigint             null,
    create_by           bigint             null,
    create_time         datetime           null,
    update_by           bigint             null,
    update_time         datetime           null,
    deleted             int    default 0   null comment '删除标志（0代表未删除，1代表已删除）'
)
    comment '评论表';

create table ancy_link
(
    id          bigint auto_increment
        primary key,
    name        varchar(256)     null,
    logo        varchar(256)     null,
    description varchar(512)     null,
    address     varchar(128)     null comment '网站地址',
    status      char default '2' null comment '审核状态 (0代表审核通过，1代表审核未通过，2代表未审核)',
    create_by   bigint           null,
    create_time datetime         null,
    update_by   bigint           null,
    update_time datetime         null,
    deleted     int  default 0   null comment '删除标志（0代表未删除，1代表已删除）'
)
    comment '友链';

create table ancy_note
(
    id          bigint auto_increment
        primary key,
    content     longtext            null comment '文章内容',
    is_top      char    default '0' null comment '是否置顶（0否，1是）',
    status      char    default '1' null comment '状态（0已发布，1草稿）',
    order_num   tinyint default 0   not null comment '排序',
    view_count  bigint  default 0   null comment '访问量',
    is_comment  char    default '1' null comment '是否允许评论 1是，0否',
    create_by   bigint              null,
    create_time datetime            null,
    update_by   bigint              null,
    update_time datetime            null,
    deleted     int     default 0   null comment '删除标志（0代表未删除，1代表已删除）'
)
    comment 'note表';

create table ancy_project
(
    id          bigint auto_increment
        primary key,
    title       varchar(256)        null comment '标题',
    content     longtext            null comment 'content',
    summary     varchar(1024)       null comment '文章摘要',
    thumbnail   varchar(256)        null comment '缩略图',
    is_top      char    default '0' null comment '是否置顶（0否，1是）',
    status      char    default '1' null comment '状态（0已发布，1草稿）',
    type        char    default '0' not null comment '0 active 1 archived',
    src_url     varchar(300)        null,
    display_url varchar(300)        null,
    order_num   tinyint default 0   not null comment '排序',
    begin_date  date                null,
    create_by   bigint              null,
    create_time datetime            null,
    update_by   bigint              null,
    update_time datetime            null,
    deleted     int     default 0   null comment '删除标志（0代表未删除，1代表已删除）'
)
    comment 'project表';

create table ancy_read
(
    id          bigint auto_increment
        primary key,
    source      varchar(256)      null comment '出处',
    content     varchar(1500)     null comment '内容',
    author      varchar(256)      null comment '作者',
    add_from    tinyint default 0 not null comment '添加来源 0 手动添加 1安读',
    create_by   bigint            null,
    create_time datetime          null,
    update_by   bigint            null,
    update_time datetime          null,
    deleted     int     default 0 null comment '删除标志（0代表未删除，1代表已删除）'
)
    comment '阅读表';

create table ancy_tag
(
    id          bigint auto_increment
        primary key,
    name        varchar(128)  null comment '标签名',
    remark      varchar(500)  null comment '备注',
    create_by   bigint        null,
    create_time datetime      null,
    update_by   bigint        null,
    update_time datetime      null,
    deleted     int default 0 null comment '删除标志（0代表未删除，1代表已删除）'
)
    comment '标签';

create table sys_menu
(
    id          bigint auto_increment comment '菜单ID'
        primary key,
    menu_name   varchar(50)              not null comment '菜单名称',
    parent_id   bigint       default 0   null comment '父菜单ID',
    order_num   int          default 0   null comment '显示顺序',
    path        varchar(200) default ''  null comment '路由地址',
    component   varchar(255)             null comment '组件路径',
    is_frame    int          default 1   null comment '是否为外链（0是 1否）',
    menu_type   char         default ''  null comment '菜单类型（M目录 C菜单 F按钮）',
    visible     char         default '0' null comment '菜单状态（0显示 1隐藏）',
    status      char         default '0' null comment '菜单状态（0正常 1停用）',
    perms       varchar(100)             null comment '权限标识',
    icon        varchar(100) default '#' null comment '菜单图标',
    create_by   bigint                   null comment '创建者',
    create_time datetime                 null comment '创建时间',
    update_by   bigint                   null comment '更新者',
    update_time datetime                 null comment '更新时间',
    remark      varchar(500) default ''  null comment '备注',
    deleted     char         default '0' null
)
    comment '菜单权限表' charset = utf8mb3;

INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1, '内容管理', 0, 1, 'content', null, 1, 'M', '0', '0', null, 'Document', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (2, '系统管理', 0, 2, 'system', null, 1, 'M', '0', '0', '', 'Setting', 1, NOW(), null, null, '系统管理目录', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (101, '写文章', 1, 1, 'write', 'content/article/write/index', 1, 'C', '0', '0', 'content:article:writer', 'EditPen', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (102, '文章管理', 1, 2, 'article', 'content/article/index', 1, 'C', '0', '0', 'content:article:list', 'Edit', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (103, '记日志', 1, 3, 'write-note', 'content/note/write/index', 1, 'C', '0', '0', 'content:note:writer', 'Calendar', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (104, '日志管理', 1, 4, 'note', 'content/note/index', 1, 'C', '0', '0', 'content:note:list', 'Calendar', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (105, '评论管理', 1, 5, 'comment', 'content/comment/index', 1, 'C', '0', '0', 'content:comment:list', 'Comment', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (106, '阅读管理', 1, 6, 'read', 'content/read/index', 1, 'C', '0', '0', 'content:read:list', 'Notebook', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (107, '项目管理', 1, 7, 'project', 'content/project/index', 1, 'C', '0', '0', 'content:project:list', 'Star', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (108, '分类管理', 1, 8, 'category', 'content/category/index', 1, 'C', '0', '0', 'content:category:list', 'Folder', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (109, '标签管理', 1, 9, 'tag', 'content/tag/index', 1, 'C', '0', '0', 'content:tag:list', 'PriceTag', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (110, '友链管理', 1, 10, 'link', 'content/link/index', 1, 'C', '0', '0', 'content:link:list', 'Link', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (111, '网站信息', 1, 11, 'setting', 'content/setting/index', 1, 'C', '0', '0', 'content:setting:list', 'Link', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (201, '用户管理', 2, 1, 'user', 'system/user/index', 1, 'C', '0', '0', 'system:user:list', 'User', 1, NOW(), null, null, '用户管理菜单', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (202, '角色管理', 2, 2, 'role', 'system/role/index', 1, 'C', '0', '0', 'system:role:list', 'UserFilled', 1, NOW(), null, null, '角色管理菜单', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (203, '菜单管理', 2, 3, 'menu', 'system/menu/index', 1, 'C', '0', '0', 'system:menu:list', 'Menu', 1, NOW(), null, null, '菜单管理菜单', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1011, '上传图片', 101, 1, '', null, 1, 'F', '0', '0', 'content:article:upload', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1021, '文章查询', 102, 1, '', null, 1, 'F', '0', '0', 'content:article:query', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1022, '文章新增', 102, 2, '', null, 1, 'F', '0', '0', 'content:article:add', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1023, '文章修改', 102, 3, '', null, 1, 'F', '0', '0', 'content:article:edit', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1024, '文章删除', 102, 4, '', null, 1, 'F', '0', '0', 'content:article:remove', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1041, '日志查询', 104, 1, '', null, 1, 'F', '0', '0', 'content:note:query', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1042, '日志新增', 104, 2, '', null, 1, 'F', '0', '0', 'content:note:add', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1043, '日志修改', 104, 3, '', null, 1, 'F', '0', '0', 'content:note:edit', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1044, '日志删除', 104, 4, '', null, 1, 'F', '0', '0', 'content:note:remove', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1051, '评论查询', 105, 1, '', null, 1, 'F', '0', '0', 'content:comment:query', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1052, '评论新增', 105, 2, '', null, 1, 'F', '0', '0', 'content:comment:add', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1053, '评论修改', 105, 3, '', null, 1, 'F', '0', '0', 'content:comment:edit', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1054, '评论删除', 105, 4, '', null, 1, 'F', '0', '0', 'content:comment:remove', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1061, '阅读查询', 106, 1, '', null, 1, 'F', '0', '0', 'content:read:query', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1062, '阅读新增', 106, 2, '', null, 1, 'F', '0', '0', 'content:read:add', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1063, '阅读修改', 106, 3, '', null, 1, 'F', '0', '0', 'content:read:edit', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1064, '阅读删除', 106, 4, '', null, 1, 'F', '0', '0', 'content:read:remove', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1071, '项目查询', 107, 1, '', null, 1, 'F', '0', '0', 'content:project:query', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1072, '项目新增', 107, 2, '', null, 1, 'F', '0', '0', 'content:project:add', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1073, '项目修改', 107, 3, '', null, 1, 'F', '0', '0', 'content:project:edit', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1074, '项目删除', 107, 4, '', null, 1, 'F', '0', '0', 'content:project:remove', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1081, '分类查询', 108, 1, '', null, 1, 'F', '0', '0', 'content:category:query', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1082, '分类新增', 108, 2, '', null, 1, 'F', '0', '0', 'content:category:add', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1083, '分类修改', 108, 3, '', null, 1, 'F', '0', '0', 'content:category:edit', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1084, '分类删除', 108, 4, '', null, 1, 'F', '0', '0', 'content:category:remove', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1085, '导出分类', 108, 5, '', null, 1, 'F', '0', '0', 'content:category:export', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1091, '标签查询', 109, 1, '', null, 1, 'F', '0', '0', 'content:tag:query', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1092, '标签新增', 109, 2, '', null, 1, 'F', '0', '0', 'content:tag:add', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1093, '标签修改', 109, 3, '', null, 1, 'F', '0', '0', 'content:tag:edit', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1094, '标签删除', 109, 4, '', null, 1, 'F', '0', '0', 'content:tag:remove', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1101, '友链查询', 110, 1, '', null, 1, 'F', '0', '0', 'content:link:query', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1102, '友链新增', 110, 2, '', null, 1, 'F', '0', '0', 'content:link:add', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1103, '友链修改', 110, 3, '', null, 1, 'F', '0', '0', 'content:link:edit', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1104, '友链删除', 110, 4, '', null, 1, 'F', '0', '0', 'content:link:remove', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (1111, '信息更新', 111, 1, '', null, 1, 'F', '0', '0', 'content:setting:edit', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (2011, '用户查询', 201, 1, '', '', 1, 'F', '0', '0', 'system:user:query', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (2012, '用户新增', 201, 2, '', '', 1, 'F', '0', '0', 'system:user:add', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (2013, '用户修改', 201, 3, '', '', 1, 'F', '0', '0', 'system:user:edit', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (2014, '用户删除', 201, 4, '', '', 1, 'F', '0', '0', 'system:user:remove', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (2015, '用户导出', 201, 5, '', '', 1, 'F', '0', '0', 'system:user:export', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (2016, '用户导入', 201, 6, '', '', 1, 'F', '0', '0', 'system:user:import', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (2017, '重置密码', 201, 7, '', '', 1, 'F', '0', '0', 'system:user:resetPwd', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (2021, '角色查询', 202, 1, '', '', 1, 'F', '0', '0', 'system:role:query', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (2022, '角色新增', 202, 2, '', '', 1, 'F', '0', '0', 'system:role:add', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (2023, '角色修改', 202, 3, '', '', 1, 'F', '0', '0', 'system:role:edit', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (2024, '角色删除', 202, 4, '', '', 1, 'F', '0', '0', 'system:role:remove', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (2025, '角色导出', 202, 5, '', '', 1, 'F', '0', '0', 'system:role:export', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (2031, '菜单查询', 203, 1, '', '', 1, 'F', '0', '0', 'system:menu:query', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (2032, '菜单新增', 203, 2, '', '', 1, 'F', '0', '0', 'system:menu:add', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (2033, '菜单修改', 203, 3, '', '', 1, 'F', '0', '0', 'system:menu:edit', '#', 1, NOW(), null, null, '', '0');
INSERT INTO ancy_blog.sys_menu (id, menu_name, parent_id, order_num, path, component, is_frame, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark, deleted) VALUES (2034, '菜单删除', 203, 4, '', '', 1, 'F', '0', '0', 'system:menu:remove', '#', 1, NOW(), null, null, '', '0');

create table sys_operate_log
(
    id            bigint auto_increment comment 'ID'
        primary key,
    operate_user  bigint        null comment '操作人ID',
    ip            varchar(100)  null comment '操作人IP',
    address       varchar(300)  null comment '操作人地址',
    ua            varchar(1000) null comment '操作人UA',
    operate_time  datetime      null comment '操作时间',
    summary       varchar(500)  null,
    class_name    varchar(100)  null comment '操作的类名',
    method_name   varchar(100)  null comment '操作的方法名',
    method_params varchar(1000) null comment '方法参数',
    return_value  varchar(2000) null comment '返回值',
    cost_time     bigint        null comment '方法执行耗时, 单位:ms'
)
    comment '操作日志表';

create table sys_role
(
    id          bigint auto_increment comment '角色ID'
        primary key,
    role_name   varchar(30)      not null comment '角色名称',
    role_key    varchar(100)     not null comment '角色权限字符串',
    role_sort   int              not null comment '显示顺序',
    status      char             not null comment '角色状态（0正常 1停用）',
    deleted     char default '0' null comment '删除标志（0代表存在 1代表删除）',
    create_by   bigint           null comment '创建者',
    create_time datetime         null comment '创建时间',
    update_by   bigint           null comment '更新者',
    update_time datetime         null comment '更新时间',
    remark      varchar(500)     null comment '备注'
)
    comment '角色信息表' charset = utf8mb3;

create table sys_role_menu
(
    id      bigint auto_increment comment 'ID'
        primary key,
    role_id bigint not null comment '角色ID',
    menu_id bigint not null comment '菜单ID'
)
    comment '角色和菜单关联表' charset = utf8mb3;

create table sys_setting
(
    id        bigint auto_increment comment 'ID'
        primary key,
    status    char         default '0'   not null comment '0正常 1停用',
    name      varchar(255) default 'key' not null,
    order_num int          default 0     not null,
    comment   varchar(255)               null,
    type      tinyint      default 0     not null comment '1基础设置 2链接 3页脚',
    value     longtext                   null,
    constraint sys_setting_pk_2
        unique (name)
)
    comment '网站设置';

INSERT INTO ancy_blog.sys_setting (id, status, name, order_num, comment, type, value) VALUES (1, '0', 'avatar', 0, 'Avatar', 1, 'https://avatars.githubusercontent.com/u/91717732');
INSERT INTO ancy_blog.sys_setting (id, status, name, order_num, comment, type, value) VALUES (2, '0', 'greeting', 0, 'Say Hi!', 1, 'Hi, Anxcye here');
INSERT INTO ancy_blog.sys_setting (id, status, name, order_num, comment, type, value) VALUES (3, '0', 'role', 0, 'What I Do', 1, 'Full Stack');
INSERT INTO ancy_blog.sys_setting (id, status, name, order_num, comment, type, value) VALUES (4, '0', 'philosophy', 0, 'My Approach', 1, 'coding with love');
INSERT INTO ancy_blog.sys_setting (id, status, name, order_num, comment, type, value) VALUES (5, '0', 'address', 0, null, 1, 'https://www.anxcye.com');
INSERT INTO ancy_blog.sys_setting (id, status, name, order_num, comment, type, value) VALUES (6, '0', 'name', 0, 'name', 1, 'Ancy');
INSERT INTO ancy_blog.sys_setting (id, status, name, order_num, comment, type, value) VALUES (7, '0', 'badge_github', 2, null, 2, '{"img":"https://github.githubassets.com/favicons/favicon.svg","index":"github","orderNum":2,"title":"Github","url":"https://github.com/Anxcye"}');
INSERT INTO ancy_blog.sys_setting (id, status, name, order_num, comment, type, value) VALUES (8, '0', 'badge_github2', 1, null, 2, '{"img":"https://github.githubassets.com/favicons/favicon.svg","index":"github2","orderNum":1,"title":"Github","url":"https://github.com/Anxcye"}');
INSERT INTO ancy_blog.sys_setting (id, status, name, order_num, comment, type, value) VALUES (9, '0', 'footer_about', 2, 'footer1', 3, '{"index":"about","orderNum":2,"position":1,"text":"about","url":"/home/9"}');
INSERT INTO ancy_blog.sys_setting (id, status, name, order_num, comment, type, value) VALUES (10, '0', 'footer_baidu', 2, 'footer2', 3, '{"index":"baidu","orderNum":2,"position":2,"text":"baidu","url":"https://baidu.com"}');
INSERT INTO ancy_blog.sys_setting (id, status, name, order_num, comment, type, value) VALUES (12, '0', 'footer_baidu2', 1, 'footer2', 3, '{"index":"baidu2","orderNum":1,"position":2,"text":"baidu2","url":"https://baidu.com"}');
INSERT INTO ancy_blog.sys_setting (id, status, name, order_num, comment, type, value) VALUES (13, '0', 'badge_aa', 3, null, 2, '{"img":"aa","index":"aa","orderNum":3,"title":"aabbc","url":"aa"}');

create table sys_user
(
    id          bigint auto_increment comment '主键'
        primary key,
    user_name   varchar(64) default 'NULL' not null comment '用户名',
    nick_name   varchar(64) default 'NULL' not null comment '昵称',
    password    varchar(64) default 'NULL' not null comment '密码',
    type        char        default '0'    null comment '用户类型：0代表普通用户，1代表管理员',
    status      char        default '0'    null comment '账号状态（0正常 1停用）',
    email       varchar(64)                null comment '邮箱',
    phonenumber varchar(32)                null comment '手机号',
    sex         char                       null comment '用户性别（0男，1女，2未知）',
    avatar      varchar(128)               null comment '头像',
    create_by   bigint                     null comment '创建人的用户id',
    create_time datetime                   null comment '创建时间',
    update_by   bigint                     null comment '更新人',
    update_time datetime                   null comment '更新时间',
    deleted     int         default 0      null comment '删除标志（0代表未删除，1代表已删除）'
)
    comment '用户表';

create table sys_user_role
(
    id      bigint auto_increment comment 'ID'
        primary key,
    user_id bigint not null comment '用户ID',
    role_id bigint not null comment '角色ID'
)
    comment '用户和角色关联表' charset = utf8mb3;


