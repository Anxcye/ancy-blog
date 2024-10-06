package com.anxcye.constants;

public class SystemConstants
{
    /**
     *  文章是草稿
     */
    public static final int ARTICLE_STATUS_DRAFT = 1;
    /**
     *  文章是正常分布状态
     */
    public static final int ARTICLE_STATUS_NORMAL = 0;

    /**
     *  状态正常
     */
    public static final String STATUS_NORMAL = "0";

    /**
     *  分类状态正常
     */
    public static final String CATEGORY_STATUS_NORMAL = STATUS_NORMAL;

    /**
     *  友链状态正常
     */
    public static final String LINK_STATUS_NORMAL = STATUS_NORMAL;

    /**
     *  评论状态正常
     */
    public static final String COMMENT_STATUS_NORMAL = STATUS_NORMAL;

    /**
     *  评论是根评论
     */
    public static final long COMMENT_IS_ROOT = -1;

    public static final String COMMENT_TYPE_ARTICLE = "0";

    public static final String COMMENT_TYPE_LINK = "1";

    public static final long SUPER_ADMIN_ID = 1L;
    public static final String MENU_TABLE_MENU = "C";
    public static final String MENU_TABLE_BUTTON = "F";
    public static final String MENU_TABLE_CATALOG = "M";
    public static final Long ROOT_MENU_PARENT_ID = 0L;
    public static final String EXPORT_CATEGORY_FILE_NAME = "分类.xlsx";
    public static final String USER_ADMIN = "1";
    public static final String USER_ENABLE = "0";
    public static final Integer ARTICLE_TYPE_NORMAL = 0;
    public static final Integer ARTICLE_TYPE_FRONT = 1;
    public static final Integer ARTICLE_TYPE_PROJECT = 2;
    public static final Integer ARTICLE_TYPE_LINK = 3;
    public static final String LINK_STATUS_HIDE = "1";
}