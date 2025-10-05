import { Team } from '@/types/player';
import type { Composer } from 'vue-i18n';

export function teamToString(team: Team, i18n: Composer) {
  switch (team) {
    case Team.Red:
      return i18n.t('teams.red');
    case Team.Blue:
      return i18n.t('teams.blue');
    case Team.Unassigned:
      return i18n.t('teams.unassigned');
    default:
      return '';
  }
}
